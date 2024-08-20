package reader

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
)

type KeyValueRepo interface {
	Get(key string) (interface{}, error)
}

type Service struct {
	KeyValueRepo KeyValueRepo
	Logger       *zap.Logger
}

func NewReaderService(repo KeyValueRepo, logger *zap.Logger) *Service {
	return &Service{
		Logger:       logger,
		KeyValueRepo: repo,
	}
}

func (s Service) Read(keys []string) (interface{}, error) {
	return s.fetchValues(keys)
}

func (s Service) fetchValues(keys []string) (map[string]interface{}, error) {
	var wg sync.WaitGroup
	results := make(map[string]interface{})
	mu := sync.Mutex{}
	errChan := make(chan error, len(keys))

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			value, err := s.KeyValueRepo.Get(k)
			if err != nil {
				wrappedErr := fmt.Errorf("error getting value for key %s: %w", k, err)
				s.Logger.Error("Error while getting value", zap.Error(wrappedErr))
				errChan <- wrappedErr
				return
			}

			v, err := typeAssert(k, value)
			if err != nil {
				s.Logger.Error("Type assertion error", zap.Error(err))
				errChan <- err
				return
			}

			mu.Lock()
			results[k] = v
			mu.Unlock()
		}(key)
	}

	wg.Wait()
	close(errChan)

	var combinedErrs []string
	for err := range errChan {
		if err != nil {
			combinedErrs = append(combinedErrs, err.Error())
		}
	}

	if len(combinedErrs) > 0 {
		return nil, errors.New(strings.Join(combinedErrs, "; "))
	}

	return results, nil
}

func typeAssert(key string, value interface{}) (any, error) {
	switch v := value.(type) {
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i, nil
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, nil
		}
		return v, nil
	case int, float64:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type for key %s: %v", key, v)
	}
}
