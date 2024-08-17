//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire
package main

import (
	"VK-test/application/auth"
	"VK-test/application/reader"
	"VK-test/application/tokens"
	"VK-test/application/writer"
	"VK-test/handler/authhandler"
	"VK-test/handler/readhandler"
	"VK-test/handler/writehandler"
	"VK-test/infrastructure/http"
	"VK-test/infrastructure/logger"
	"VK-test/infrastructure/repo/kv_storage"
	"VK-test/infrastructure/repo/user_data"
	tarantoolclient "VK-test/infrastructure/tarantool"
	"github.com/google/wire"
)

// SuperSet is the Wire provider set that includes all the dependencies.
var SuperSet = wire.NewSet(
	kv_storage.NewTarantoolDataRepo,
	tokens.NewPayloadService,
	user_data.NewTarantoolUserRepo,
	auth.NewService,
	readhandler.NewHandler,
	authhandler.NewHandler,
	writehandler.NewHandler,
	reader.NewReaderService,
	writer.NewWriterService,
	http.NewHandlerContainer,
	logger.NewLogger,
	tarantoolclient.NewTarantoolClient,
	NewApp,

	// Bind interfaces to their implementations
	wire.Bind(new(writer.KeyValueRepo), new(*kv_storage.TarantoolDataRepo)),
	wire.Bind(new(auth.TokenService), new(*tokens.PayloadService)),
	wire.Bind(new(auth.UserRepo), new(*user_data.TarantoolUserRepo)),
	wire.Bind(new(readhandler.ReaderService), new(*reader.Service)),
	wire.Bind(new(writehandler.WriterService), new(*writer.Service)),
	wire.Bind(new(http.AuthHandler), new(*authhandler.Handler)),
	wire.Bind(new(http.ReadHandler), new(*readhandler.Handler)),
	wire.Bind(new(http.WriteHandler), new(*writehandler.Handler)),
	wire.Bind(new(authhandler.ServiceAuth), new(*auth.Service)),
	wire.Bind(new(reader.KeyValueRepo), new(*kv_storage.TarantoolDataRepo)),
)

func InitializeApp() (*App, error) {
	wire.Build(SuperSet)
	return &App{}, nil
}
