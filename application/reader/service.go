package reader

import "sync"

type KeyValueRepo interface {
	Get(key string) (interface{}, error)
}

type Service struct {
	KeyValueRepo KeyValueRepo
}

func NewReaderService(repo KeyValueRepo) *Service {
	return &Service{
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
				errChan <- err
				return
			}
			mu.Lock()
			results[k] = value
			mu.Unlock()
		}(key)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
