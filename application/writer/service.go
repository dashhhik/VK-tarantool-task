package writer

import (
	"fmt"
	"strings"
	"sync"
)

type KeyValueRepo interface {
	Set(key string, value interface{}) error
}

type Service struct {
	UserRepo KeyValueRepo
}

func NewWriterService(repo KeyValueRepo) *Service {
	return &Service{
		UserRepo: repo,
	}
}

func (s Service) Write(data map[string]interface{}) error {
	return s.processMap(data)
}

func (s Service) processMap(data map[string]interface{}) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(data))

	for key, value := range data {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			if err := s.UserRepo.Set(k, v); err != nil {
				mu.Lock()
				errChan <- fmt.Errorf("failed to set key %s: %w", k, err)
				mu.Unlock()
			}
		}(key, value)
	}

	wg.Wait()
	close(errChan)

	var combinedErrors []string
	for err := range errChan {
		if err != nil {
			combinedErrors = append(combinedErrors, err.Error())
		}
	}

	if len(combinedErrors) > 0 {
		return fmt.Errorf("multiple errors occurred: %s", strings.Join(combinedErrors, "; "))
	}

	return nil
}
