package executer

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	story     *Executor
	requester *MockRequester
	limiter   *MockLimiter
	ctl       *gomock.Controller
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.requester = NewMockRequester(s.ctl)
	s.limiter = NewMockLimiter(s.ctl)

}

func (s *Suite) Test_01() {
	s.T().Log(`Настройки:
	            - 10 параллельных запросов к внешним сервисам
			    - 1  паралельный запрос к нашему сервису
			    - 2  URL
			   Запросы:
			    - 1 паралельный запрос`)

	yaUrl, _ := url.Parse("yandex.ru")
	googleUrl, _ := url.Parse("google.com")
	ctx := context.Background()

	s.limiter.EXPECT().Take().Times(1).Return(make(chan bool, 1))
	s.limiter.EXPECT().Free().Times(1)
	config := Config{RequestsInParallel: 10, RequestTimeout: time.Second, MaxUrlNum: 20}
	s.story = New(config, s.limiter, s.requester)

	s.requester.EXPECT().Do(gomock.Any(), yaUrl).Return("use google", nil)
	s.requester.EXPECT().Do(gomock.Any(), googleUrl).Return("use yandex", nil)

	resp, err := s.story.MakeRequest(ctx, []string{"yandex.ru", "google.com"})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(resp))
}

func (s *Suite) Test_02() {
	s.T().Log(`Настройки:
	            - 1 параллельый запросов к внешним сервисам
			    - 1 паралельный запрос к нашему сервису
				- 2  URL
			   Запросы:
			    - 1 паралельный запрос`)

	yaUrl, _ := url.Parse("yandex.ru")
	googleUrl, _ := url.Parse("google.com")
	ctx := context.Background()

	s.limiter.EXPECT().Take().Times(1).Return(make(chan bool, 1))
	s.limiter.EXPECT().Free().Times(1)
	config := Config{RequestsInParallel: 1, RequestTimeout: time.Second, MaxUrlNum: 20}
	s.story = New(config, s.limiter, s.requester)

	s.requester.EXPECT().Do(gomock.Any(), yaUrl).Return("use google", nil)
	s.requester.EXPECT().Do(gomock.Any(), googleUrl).Return("use yandex", nil)

	resp, err := s.story.MakeRequest(ctx, []string{"yandex.ru", "google.com"})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(resp))
}

func (s *Suite) Test_03() {
	s.T().Log(`Настройки:
	            - 1 параллельый запросов к внешним сервисам
			    - 0 паралельный запрос к нашему сервису
			   Запросы:
			    - 1 паралельный запрос
			   Результат:
			    - запрос прерван по таймауту`)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	s.limiter.EXPECT().Take().Times(1).Return(make(chan bool))
	config := Config{RequestsInParallel: 1, RequestTimeout: time.Second, MaxUrlNum: 20}
	s.story = New(config, s.limiter, s.requester)

	_, err := s.story.MakeRequest(ctx, []string{"yandex.ru", "google.com"})
	require.True(s.T(), errors.Is(err, context.DeadlineExceeded))
}

func (s *Suite) Test_04() {
	s.T().Log(`Настройки:
	            - 1 параллельный запрос к внешним сервисам
			    - 1 паралельный запрос к нашему сервису
				- 2  URL
				- внешений сервис возвращает ошибку
			   Запросы:
			    - 1 паралельный запрос`)

	ctx := context.Background()

	s.limiter.EXPECT().Take().Times(1).Return(make(chan bool, 1))
	s.limiter.EXPECT().Free().Times(1)
	config := Config{RequestsInParallel: 1, RequestTimeout: time.Second, MaxUrlNum: 20}
	s.story = New(config, s.limiter, s.requester)

	s.requester.EXPECT().Do(gomock.Any(), gomock.Any()).AnyTimes().Return("use google", fmt.Errorf("some error"))

	_, err := s.story.MakeRequest(ctx, []string{"yandex.ru", "google.com"})
	require.Error(s.T(), err)
}

func (s *Suite) Test_05() {
	s.T().Log(`Настройки:
	            - 1 параллельный запрос к внешним сервисам
			    - 1 паралельный запрос к нашему сервису
				- 1  URL
			   Запросы:
			    - 1 паралельный запрос`)

	yaUrl, _ := url.Parse("yandex.ru")
	ctx := context.Background()

	s.limiter.EXPECT().Take().Times(1).Return(make(chan bool, 1))
	s.limiter.EXPECT().Free().Times(1)
	config := Config{RequestsInParallel: 1, RequestTimeout: time.Second, MaxUrlNum: 20}
	s.story = New(config, s.limiter, s.requester)

	s.requester.EXPECT().Do(gomock.Any(), yaUrl).Return("use google", fmt.Errorf("some error"))

	_, err := s.story.MakeRequest(ctx, []string{"yandex.ru"})
	require.Error(s.T(), err)
}

func (s *Suite) TearDownTest() {
	s.ctl.Finish()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
