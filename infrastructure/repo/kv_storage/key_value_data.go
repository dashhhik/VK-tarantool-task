package kv_storage

import (
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
)

const keyValueSpace = "key_value"

type KeyValue struct {
	_msgpack struct{} `msgpack:",asArray"`
	Key      string
	Value    interface{}
}

type TarantoolDataRepo struct {
	Client *tarantool.Connection
}

func NewTarantoolDataRepo(client *tarantool.Connection) *TarantoolDataRepo {
	return &TarantoolDataRepo{
		Client: client,
	}
}

func (r TarantoolDataRepo) Get(key string) (interface{}, error) {
	var value string

	req := tarantool.NewSelectRequest(keyValueSpace).Key(key).Limit(1)
	err := r.Client.Do(req).GetTyped(&value)
	if err != nil {
		fmt.Printf("error while getting value: %v", err)
		return "", err
	}

	return value, nil
}

func (r TarantoolDataRepo) Set(key string, value interface{}) error {
	kv := KeyValue{
		Key:   key,
		Value: value,
	}
	_, err := r.Client.Do(tarantool.NewReplaceRequest(keyValueSpace).Tuple(kv)).Get()
	if err != nil {
		return fmt.Errorf("failed to insert key-value pair: %w", err)
	}

	return nil
}
