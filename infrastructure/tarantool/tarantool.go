package tarantool

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
)

func NewTarantoolClient() *tarantool.Connection {
	ctx := context.Background()

	dialer := tarantool.NetDialer{
		User:     "admin",
		Password: "presale",
		Address:  "127.0.0.1:3301",
	}

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})
	if err != nil {
		fmt.Println("Connection refused:", err)
		return nil
	}
	return conn

}
