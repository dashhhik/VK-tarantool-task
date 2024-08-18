package tarantool

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"time"
)

func NewTarantoolClient() *tarantool.Connection {
	time.Sleep(3 * time.Second)
	dialer := tarantool.NetDialer{
		User:     "admin",
		Password: "presale",
		Address:  "tarantool:3301",
	}

	opts := tarantool.Opts{
		Timeout:       5 * time.Second,
		Concurrency:   32,
		Reconnect:     time.Second,
		MaxReconnects: 10,
	}

	var err error
	var conn *tarantool.Connection

	for i := uint(0); i < opts.MaxReconnects; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		conn, err = tarantool.Connect(ctx, dialer, opts)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(opts.Reconnect)
	}
	if err != nil {
		fmt.Println(err, "Failed to connect to Tarantool")
		return nil
	}

	data, err := conn.Do(tarantool.NewPingRequest()).Get()
	fmt.Println("Ping Data", data)
	fmt.Println("Ping Error", err)

	fmt.Println("Connection is ready and ping successful")
	return conn
}
