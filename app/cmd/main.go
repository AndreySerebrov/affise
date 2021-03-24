package main

import (
	"affice/internal/common/env"
	"affice/internal/domains/requests/gates/requester"
	v1 "affice/internal/domains/requests/handlers/http/v1"
	"affice/internal/domains/requests/stories/executer"
	"affice/internal/domains/requests/stories/loclimiter"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	envVals := env.LoadEnv()

	logger := log.New(os.Stdout, "", 0)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	gate := requester.New()
	inputLim := loclimiter.New(envVals.InputRateLimit)

	config := executer.Config{
		RequestsInParallel: envVals.OutputRateLimit,
		RequestTimeout:     envVals.RequestTimeout,
		MaxUrlNum:          envVals.MaxUrlNum}

	ex := executer.New(config, inputLim, gate)

	log.Println("")
	log.Println("Параметры приложения:")
	log.Println("")
	log.Printf("Таймаут запроса к сервису: %s \n", config.RequestTimeout.String())
	log.Printf("Количество URL в запросе не более: %d \n", config.MaxUrlNum)
	log.Printf("Запросов принимается параллельно не более: %d \n", envVals.InputRateLimit)
	log.Printf("Исходящих запросов на один входящий единовременно не более: %d \n", envVals.OutputRateLimit)
	log.Println("")
	http.HandleFunc("/", v1.New(ctx, logger, ex).Handle)

	go func() {
		log.Printf("Сервер стартовал на порту %d", envVals.AppPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", envVals.AppPort), nil)
		if err != nil {
			log.Println(err.Error())
			cancel()
		}
	}()

	oscall := <-c
	log.Printf("Завершение работы")
	log.Printf("system call:%+v", oscall)
	cancel()
}
