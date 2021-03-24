package executer

import (
	"affice/internal/domains/requests/stories/loclimiter"
	"context"
	"fmt"
	"net/url"
	"sync"
)

type Executor struct {
	limiter   Limiter
	requester Requester
	config    Config
}

func New(config Config, limiter Limiter, requester Requester) *Executor {
	return &Executor{
		limiter:   limiter,
		requester: requester,
		config:    config,
	}
}

func (e *Executor) MakeRequest(ctx context.Context, urls []string) (responseList []Response, errFun error) {

	if len(urls) > int(e.config.MaxUrlNum) {
		return nil, fmt.Errorf("Url number should be equal or less than %d", e.config.MaxUrlNum)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e.limiter.Take() <- true:
	}
	defer e.limiter.Free()

	locLimiter := loclimiter.New(e.config.RequestsInParallel)
	errChan := make(chan error)
	responseChan := make(chan Response)
	taskChan := make(chan string)
	doneChan := make(chan bool, 1)
	wg := &sync.WaitGroup{}

	// Рутина выполнения запросов
	go func() {
		for task := range taskChan {

			go func(rawurl string) {
				defer wg.Done()

				parsedURL, err := url.Parse(rawurl)
				if err != nil {
					errChan <- err
					return
				}
				ctxTimeout, cancel := context.WithTimeout(ctx, e.config.RequestTimeout)
				defer cancel()
				response, err := e.requester.Do(ctxTimeout, parsedURL)
				if err != nil {
					errChan <- err
					return
				}
				responseChan <- Response{Url: rawurl, Response: response}
				locLimiter.Free()
			}(task)
		}
	}()

	// Здесь ловим ошибки
	go func() {
		select {
		case errFun = <-errChan:
		case <-ctx.Done():
			errFun = ctx.Err()
		}
		close(doneChan)
	}()

	// Постановка задач
	go func() {
		for _, rawurl := range urls {
			select {
			case <-doneChan:
				break
			case locLimiter.Take() <- true:
				wg.Add(1)
				taskChan <- rawurl
			}
		}

		wg.Wait()
		close(responseChan)
	}()

	// Получение результата
	for response := range responseChan {
		responseList = append(responseList, response)
	}

	close(taskChan)
	close(errChan)

	return responseList, errFun
}
