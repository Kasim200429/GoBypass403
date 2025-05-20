package bypass

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RunAllBypassTechniques executes all bypass techniques against the target URL
func RunAllBypassTechniques(config Config) ([]Result, error) {
	var results []Result
	var wg sync.WaitGroup
	var mutex sync.Mutex
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	techniques := GetTechniques()
	resultCh := make(chan []Result, len(techniques))
	errCh := make(chan error, len(techniques))

	for _, technique := range techniques {
		wg.Add(1)
		go func(t Technique) {
			defer wg.Done()

			techniqueResults, err := t.Test(config.URL, client, config)
			if err != nil {
				errCh <- fmt.Errorf("error in %s: %v", t.Name, err)
				return
			}

			if len(techniqueResults) > 0 {
				resultCh <- techniqueResults
			}
		}(technique)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	for res := range resultCh {
		mutex.Lock()
		results = append(results, res...)
		mutex.Unlock()
	}

	// Check for errors
	for err := range errCh {
		if err != nil {
			return results, err
		}
	}

	return results, nil
}
