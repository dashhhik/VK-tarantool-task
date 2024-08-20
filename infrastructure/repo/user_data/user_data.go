package user_data

import (
	"VK-test/core"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

const userSpace = "users"

type TarantoolUserRepo struct {
	Client *tarantool.Connection
	Logger *zap.Logger
}

func NewTarantoolUserRepo(client *tarantool.Connection, logger *zap.Logger) *TarantoolUserRepo {
	return &TarantoolUserRepo{
		Client: client,
		Logger: logger,
	}
}

func (r *TarantoolUserRepo) Get(key string) (string, error) {
	req := tarantool.NewSelectRequest(userSpace).Key([]interface{}{key}).Limit(1)
	resp, err := r.Client.Do(req).GetResponse()
	if err != nil {
		r.Logger.Error("Error while getting value from Tarantool", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	decoded, err := resp.Decode()
	if err != nil {
		r.Logger.Error("Error decoding response from Tarantool", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("failed to decode response for key %s: %w", key, err)
	}

	r.Logger.Info("Decoded response", zap.Any("decoded", decoded))

	if len(decoded) == 0 {
		err := fmt.Errorf("%w: %s", core.ErrKeyNotFound, key)
		r.Logger.Warn("No result found for key", zap.String("key", key), zap.Error(err))
		return "", err
	}

	row := decoded[0].([]interface{})
	password, ok := row[1].(string)
	if !ok {
		err := fmt.Errorf("unexpected type for password field for key %s", key)
		r.Logger.Error("Error processing password field", zap.Error(err))
		return "", err
	}

	r.Logger.Info("Successfully retrieved password", zap.String("key", key))
	return password, nil
}
