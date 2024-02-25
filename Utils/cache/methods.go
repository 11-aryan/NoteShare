package cache

import (
	"context"
	"encoding/json"
)

func (r *RedisClient) SetHash(key string, field string, data interface{}) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = r.client.HSet(context.TODO(), key, field, string(jsonData)).Err()
	if err != nil {
		return
	}
	return
}

func (r *RedisClient) GetHash(key string, field string) (data map[string][]string, err error) {
	jsonData, err := r.client.HGet(context.TODO(), key, field).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return
	}
	return
}
