package kv_storage

import (
	"errors"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

const keyValueSpace = "key_value"

type KeyValue struct {
	_msgpack struct{} `msgpack:",asArray"`
	Key      string
	Value    interface{}
}

type TarantoolDataRepo struct {
	Client *tarantool.Connection
	Logger *zap.Logger
}

func NewTarantoolDataRepo(client *tarantool.Connection, logger *zap.Logger) *TarantoolDataRepo {
	return &TarantoolDataRepo{
		Client: client,
		Logger: logger,
	}
}

var ErrKeyNotFound = errors.New("key not found")

func (r *TarantoolDataRepo) Get(key string) (interface{}, error) {
	req := tarantool.NewSelectRequest(keyValueSpace).Key([]interface{}{key}).Limit(1)
	resp, err := r.Client.Do(req).GetResponse()
	if err != nil {
		r.Logger.Error("Error while getting value from Tarantool", zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	decoded, err := resp.Decode()
	if err != nil {
		r.Logger.Error("Error decoding response from Tarantool", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response for key %s: %w", key, err)
	}

	if len(decoded) == 0 {
		err := fmt.Errorf("%w: %s", ErrKeyNotFound, key)
		r.Logger.Warn("No result found for key", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	row := decoded[0].([]interface{})
	value := row[1]

	r.Logger.Info("Successfully retrieved value", zap.String("key", key), zap.Any("value", value))
	return value, nil
}

func (r *TarantoolDataRepo) Set(key string, value interface{}) error {
	kv := KeyValue{
		Key:   key,
		Value: value,
	}
	_, err := r.Client.Do(tarantool.NewReplaceRequest(keyValueSpace).Tuple(kv)).Get()
	if err != nil {
		r.Logger.Error("Error inserting key-value pair into Tarantool", zap.String("key", key), zap.Any("value", value), zap.Error(err))
		return fmt.Errorf("failed to insert key-value pair %s: %w", key, err)
	}

	r.Logger.Info("Successfully inserted key-value pair", zap.String("key", key), zap.Any("value", value))
	return nil
}
