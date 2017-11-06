// Package storage provides a type for storing Go objects in Redis.
package storage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// ErrNotFound is the error returned when the given key is not found.
var ErrNotFound = errors.New("not found")

// Storage is the basic type that provides methods for saving, listing and
// deleting types on Redis.
type Storage struct {
	once   sync.Once
	config *Config
	client *redis.Client
}

// Config contains configuration for the Redis, in the standard proposed by
// Gizmo.
type Config struct {
	// Comma-separated list of sentinel servers.
	//
	// Example: 10.10.10.10:6379,10.10.10.1:6379,10.10.10.2:6379.
	SentinelAddrs      string `envconfig:"SENTINEL_ADDRS"`
	SentinelMasterName string `envconfig:"SENTINEL_MASTER_NAME"`

	RedisAddr          string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
	Password           string `envconfig:"REDIS_PASSWORD"`
	PoolSize           int    `envconfig:"REDIS_POOL_SIZE"`
	PoolTimeout        int    `envconfig:"REDIS_POOL_TIMEOUT_SECONDS"`
	IdleTimeout        int    `envconfig:"REDIS_IDLE_TIMEOUT_SECONDS"`
	IdleCheckFrequency int    `envconfig:"REDIS_IDLE_CHECK_FREQUENCY_SECONDS"`
}

// RedisClient creates a new instance of the client using the underlying
// configuration.
func (c *Config) RedisClient() *redis.Client {
	if c.SentinelAddrs != "" {
		sentinelAddrs := strings.Split(c.SentinelAddrs, ",")
		return redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:      sentinelAddrs,
			MasterName:         c.SentinelMasterName,
			Password:           c.Password,
			PoolSize:           c.PoolSize,
			PoolTimeout:        time.Duration(c.PoolTimeout) * time.Second,
			IdleTimeout:        time.Duration(c.IdleTimeout) * time.Second,
			IdleCheckFrequency: time.Duration(c.IdleCheckFrequency) * time.Second,
		})
	}
	redisAddr := c.RedisAddr
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr:               redisAddr,
		Password:           c.Password,
		PoolSize:           c.PoolSize,
		PoolTimeout:        time.Duration(c.PoolTimeout) * time.Second,
		IdleTimeout:        time.Duration(c.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(c.IdleCheckFrequency) * time.Second,
	})
}

// NewStorage returns a new instance of storage with the given configuration.
func NewStorage(cfg *Config) (*Storage, error) {
	return &Storage{config: cfg}, nil
}

// Save creates the given key as a Redis hash.
//
// The given hash must be either a struct or map[string]string.
func (s *Storage) Save(key string, hash interface{}) error {
	fields, err := s.FieldMap(hash)
	if err != nil {
		return err
	}
	return s.RedisClient().HMSet(key, fields).Err()
}

// FieldMap extract the map of fields from the given type (which can be a
// struct, a map[string]string or pointer to those).
func (s *Storage) FieldMap(hash interface{}) (map[string]interface{}, error) {
	if hash == nil {
		return nil, errors.New("no fields provided")
	}
	value := reflect.ValueOf(hash)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map:
		return s.mapToFieldList(hash)
	case reflect.Struct:
		return s.structToFieldList(value)
	default:
		return nil, errors.New("please provide a map or a struct")
	}
}

func (s *Storage) mapToFieldList(hash interface{}, prefixes ...string) (map[string]interface{}, error) {
	m, ok := hash.(map[string]string)
	if !ok {
		return nil, errors.New("please provide a map[string]string")
	}
	if len(m) < 1 {
		return nil, errors.New("please provide a map[string]string with at least one item")
	}
	fields := make(map[string]interface{}, len(m))
	for key, value := range m {
		key = strings.Join(append(prefixes, key), "_")
		fields[key] = value
	}
	return fields, nil
}

func (s *Storage) structToFieldList(value reflect.Value, prefixes ...string) (map[string]interface{}, error) {
	fields := make(map[string]interface{})
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
				expandedFields, err := s.structToFieldList(fieldValue, myPrefixes...)
				if err != nil {
					return nil, err
				}
				for k, v := range expandedFields {
					fields[k] = v
				}
			case reflect.Map:
				expandedFields, err := s.mapToFieldList(fieldValue.Interface(), myPrefixes...)
				if err != nil {
					return nil, err
				}
				for k, v := range expandedFields {
					fields[k] = v
				}
			default:
				return nil, errors.New("can only expand structs and maps")
			}
		} else {
			if parts[0] != "" {
				key := strings.Join(append(prefixes, parts[0]), "_")
				var strValue string
				iface := fieldValue.Interface()
				switch v := iface.(type) {
				case time.Time:
					strValue = v.Format(time.RFC3339Nano)
				case []string:
					strValue = strings.Join(v, "%%%")
				default:
					strValue = fmt.Sprintf("%v", v)
				}
				if parts[len(parts)-1] == "omitempty" && strValue == "" {
					continue
				}
				fields[key] = strValue
			}
		}
	}
	return fields, nil
}

// Load loads the given key in the given output. The output must be a pointer
// to a struct or a map[string]string.
func (s *Storage) Load(key string, out interface{}) error {
	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Ptr {
		return errors.New("please provide a pointer for getting result from the database")
	}
	value = value.Elem()
	result, err := s.RedisClient().HGetAll(key).Result()
	if err != nil {
		return err
	}
	if len(result) < 1 {
		return ErrNotFound
	}
	switch value.Kind() {
	case reflect.Map:
		return s.loadMap(result, value)
	case reflect.Struct:
		return s.loadStruct(result, value)
	default:
		return errors.New("please provider a pointer to a struct or a map for getting result from the database")
	}
}

func (s *Storage) loadMap(in map[string]string, out reflect.Value, prefixes ...string) error {
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

func (s *Storage) loadStruct(in map[string]string, out reflect.Value, prefixes ...string) error {
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
				err := s.loadMap(in, fieldValue, myPrefixes...)
				if err != nil {
					return err
				}
			case reflect.Struct:
				err := s.loadStruct(in, fieldValue, myPrefixes...)
				if err != nil {
					return err
				}
			default:
				return errors.New("can only expand values to structs or maps")
			}
		} else {
			key := strings.Join(append(prefixes, parts[0]), "_")
			if value, ok := in[key]; ok {
				switch fieldValue.Kind() {
				case reflect.Slice:
					values := strings.Split(value, "%%%")
					if reflect.TypeOf(values).AssignableTo(fieldValue.Type()) {
						fieldValue.Set(reflect.ValueOf(values))
					}
				case reflect.Bool:
					boolValue, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					fieldValue.SetBool(boolValue)
				case reflect.Float64:
					floatValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return err
					}
					fieldValue.SetFloat(floatValue)
				case reflect.Int:
					intValue, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return err
					}
					fieldValue.SetInt(intValue)
				case reflect.Uint:
					uintValue, err := strconv.ParseUint(value, 10, 64)
					if err != nil {
						return err
					}
					fieldValue.SetUint(uintValue)
				case reflect.Struct:
					if reflect.TypeOf(time.Time{}).AssignableTo(fieldValue.Type()) {
						timeValue, err := time.Parse(time.RFC3339Nano, value)
						if err != nil {
							return err
						}
						fieldValue.Set(reflect.ValueOf(timeValue))
					}
				default:
					fieldValue.SetString(value)
				}
			}
		}
	}
	return nil
}

// Delete deletes the given key from redis, returning ErrNotFound when it
// doesn't exist.
func (s *Storage) Delete(key string) error {
	n, err := s.RedisClient().Del(key).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// RedisClient returns the underlying Redis client.
func (s *Storage) RedisClient() *redis.Client {
	s.once.Do(func() {
		s.client = s.config.RedisClient()
	})
	return s.client
}
