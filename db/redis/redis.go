package redis

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/db"
	"gopkg.in/redis.v3"
)

var errNotFound = errors.New("not found")

func init() {
	redis.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
}

// NewRepository creates a new Repository that uses Redis for persistence.
func NewRepository(cfg *config.Config) (db.Repository, error) {
	repo := &redisRepository{config: cfg}
	repo.client = repo.redisClient()
	return &redisRepository{config: cfg}, nil
}

type redisRepository struct {
	config *config.Config
	client *redis.Client
	once   sync.Once
}

func (r *redisRepository) save(key string, hash interface{}) error {
	fields, err := r.fieldList(hash)
	if err != nil {
		return err
	}
	return r.redisClient().HMSet(key, fields[0], fields[1], fields[2:]...).Err()
}

func (r *redisRepository) fieldList(hash interface{}) ([]string, error) {
	if hash == nil {
		return nil, errors.New("no fields provided")
	}
	value := reflect.ValueOf(hash)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map:
		return r.mapToFieldList(hash)
	case reflect.Struct:
		return r.structToFieldList(value)
	default:
		return nil, errors.New("please provide a map or a struct")
	}
}

func (r *redisRepository) mapToFieldList(hash interface{}, prefixes ...string) ([]string, error) {
	m, ok := hash.(map[string]string)
	if !ok {
		return nil, errors.New("please provide a map[string]string")
	}
	if len(m) < 1 {
		return nil, errors.New("please provide a map[string]string with at least one item")
	}
	fields := make([]string, 0, len(m)*2)
	for key, value := range m {
		key = strings.Join(append(prefixes, key), "_")
		fields = append(fields, key, value)
	}
	return fields, nil
}

func (r *redisRepository) structToFieldList(value reflect.Value, prefixes ...string) ([]string, error) {
	fields := make([]string, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldName := field.Tag.Get("redis-hash")
		if fieldName == "-" {
			continue
		}
		parts := strings.Split(fieldName, ",")
		fieldValue := value.Field(i)
		if len(parts) > 1 && parts[len(parts)-1] == "expand" {
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			myPrefixes := append(prefixes, parts[0])
			switch fieldValue.Kind() {
			case reflect.Struct:
				expandedFields, err := r.structToFieldList(fieldValue, myPrefixes...)
				if err != nil {
					return nil, err
				}
				fields = append(fields, expandedFields...)
			case reflect.Map:
				expandedFields, err := r.mapToFieldList(fieldValue.Interface(), myPrefixes...)
				if err != nil {
					return nil, err
				}
				fields = append(fields, expandedFields...)
			default:
				return nil, errors.New("can only expand structs and maps")
			}
		} else {
			if parts[0] != "" {
				key := strings.Join(append(prefixes, parts[0]), "_")
				fields = append(fields, key, fmt.Sprintf("%v", fieldValue.Interface()))
			}
		}
	}
	return fields, nil
}

func (r *redisRepository) load(key string, out interface{}) error {
	result, err := r.redisClient().HGetAllMap(key).Result()
	if err != nil {
		return err
	}
	if len(result) < 1 {
		return errNotFound
	}
	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Ptr {
		return errors.New("please provide a pointer for getting result from the database")
	}
	value = value.Elem()
	switch value.Kind() {
	case reflect.Map:
		return r.loadMap(result, value)
	case reflect.Struct:
		return r.loadStruct(result, value)
	default:
		return errors.New("please provider a pointer to a struct or a map for getting result from the database")
	}
}

func (r *redisRepository) loadMap(in map[string]string, out reflect.Value, prefixes ...string) error {
	if out.Type().Key().Kind() != reflect.String || out.Type().Elem().Kind() != reflect.String {
		return errors.New("please provide a map[string]string")
	}
	joinedPrefixes := strings.Join(prefixes, "_")
	if joinedPrefixes != "" {
		joinedPrefixes += "_"
	}
	for k, v := range in {
		if !strings.HasPrefix(k, joinedPrefixes) {
			continue
		}
		k = strings.Replace(k, joinedPrefixes, "", 1)
		out.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
	return nil
}

func (r *redisRepository) loadStruct(in map[string]string, out reflect.Value, prefixes ...string) error {
	for i := 0; i < out.NumField(); i++ {
		field := out.Type().Field(i)
		if field.PkgPath != "" {
			continue
		}
		tagValue := field.Tag.Get("redis-hash")
		if tagValue == "-" {
			continue
		}
		parts := strings.Split(tagValue, ",")
		fieldValue := out.Field(i)
		if len(parts) > 1 && parts[len(parts)-1] == "expand" {
			myPrefixes := append(prefixes, parts[0])
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			switch fieldValue.Kind() {
			case reflect.Map:
				err := r.loadMap(in, fieldValue, myPrefixes...)
				if err != nil {
					return err
				}
			case reflect.Struct:
				err := r.loadStruct(in, fieldValue, myPrefixes...)
				if err != nil {
					return err
				}
			default:
				return errors.New("can only expand values to structs or maps")
			}
		} else {
			key := strings.Join(append(prefixes, parts[0]), "_")
			if value, ok := in[key]; ok {
				if fieldValue.Kind() == reflect.Bool {
					boolValue, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					fieldValue.SetBool(boolValue)
				} else {
					fieldValue.SetString(value)
				}
			}
		}
	}
	return nil
}

func (r *redisRepository) delete(key string, notFoundErr error) error {
	n, err := r.redisClient().Del(key).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return notFoundErr
	}
	return nil
}

func (r *redisRepository) generateID() (string, error) {
	var raw [8]byte
	n, err := rand.Read(raw[:])
	if err != nil {
		return "", err
	}
	if n != 8 {
		return "", io.ErrShortWrite
	}
	return fmt.Sprintf("%x", raw), nil
}

func (r *redisRepository) redisClient() *redis.Client {
	r.once.Do(func() {
		var sentinelAddrs []string
		if r.config.Redis.SentinelAddrs != "" {
			sentinelAddrs = strings.Split(r.config.Redis.SentinelAddrs, ",")
		}
		if len(sentinelAddrs) > 0 {
			r.client = redis.NewFailoverClient(&redis.FailoverOptions{
				SentinelAddrs: sentinelAddrs,
				MasterName:    r.config.Redis.SentinelMasterName,
				Password:      r.config.Redis.Password,
				PoolSize:      r.config.Redis.PoolSize,
				PoolTimeout:   time.Duration(r.config.Redis.PoolTimeout) * time.Second,
			})
		} else {
			redisAddr := r.config.Redis.RedisAddr
			if redisAddr == "" {
				redisAddr = "127.0.0.1:6379"
			}
			r.client = redis.NewClient(&redis.Options{
				Addr:        redisAddr,
				Password:    r.config.Redis.Password,
				PoolSize:    r.config.Redis.PoolSize,
				PoolTimeout: time.Duration(r.config.Redis.PoolTimeout) * time.Second,
			})
		}
	})
	return r.client
}
