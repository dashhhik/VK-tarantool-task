package user_data

import (
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
)

const userSpace = "users"

type User struct {
	_msgpack struct{} `msgpack:",asArray"`
	Username string
	Password string
}

type TarantoolUserRepo struct {
	Client *tarantool.Connection
}

func NewTarantoolUserRepo(client *tarantool.Connection) *TarantoolUserRepo {
	return &TarantoolUserRepo{
		Client: client,
	}
}

func (r TarantoolUserRepo) Get(key string) (interface{}, error) {
	var value string

	req := tarantool.NewSelectRequest(userSpace).Key(key).Limit(1)
	err := r.Client.Do(req).GetTyped(&value)
	if err != nil {
		fmt.Printf("error while getting value: %v", err)
		return "", err
	}

	return value, nil
}
