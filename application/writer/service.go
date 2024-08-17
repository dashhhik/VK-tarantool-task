package writer

import "sync"

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
	errChan := make(chan error, len(data))

	for key, value := range data {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			if err := s.UserRepo.Set(k, v); err != nil {
				errChan <- err
			}
		}(key, value)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
