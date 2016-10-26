package redis

import "gopkg.in/redis.v4"

func cleanRedis() error {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	defer client.Close()
	err := deleteKeys("job:*", client)
	if err != nil {
		return err
	}
	err = deleteKeys("presetmap:*", client)
	if err != nil {
		return err
	}
	err = deleteKeys(presetmapsSetKey, client)
	if err != nil {
		return err
	}
	return deleteKeys(jobsSetKey, client)
}

func deleteKeys(pattern string, client *redis.Client) error {
	keys, err := client.Keys(pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = client.Del(keys...).Result()
	}
	return err
}
