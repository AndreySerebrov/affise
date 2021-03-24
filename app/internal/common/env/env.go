package env

import (
	"os"
	"strconv"
	"time"
)

type EnvVals struct {
	AppPort         int64
	InputRateLimit  int64
	OutputRateLimit int64
	RequestTimeout  time.Duration
	MaxUrlNum       int64
}

func LoadEnv() (envVals EnvVals) {
	envVals.AppPort = 8081
	envVals.InputRateLimit = 100
	envVals.OutputRateLimit = 4
	envVals.RequestTimeout = time.Second
	envVals.MaxUrlNum = 20

	appPort := os.Getenv("APPLICATION_PORT")
	if appPort != "" {
		port, err := strconv.ParseInt(appPort, 10, 64)
		if err == nil {
			envVals.AppPort = port
		}
	}

	inputRateLimit := os.Getenv("IN_RATE_LIMIT")
	if inputRateLimit != "" {
		lim, err := strconv.ParseInt(inputRateLimit, 10, 64)
		if err == nil {
			envVals.InputRateLimit = lim
		}
	}

	outRateLimit := os.Getenv("OUT_RATE_LIMIT")
	if outRateLimit != "" {
		lim, err := strconv.ParseInt(outRateLimit, 10, 64)
		if err == nil {
			envVals.OutputRateLimit = lim
		}
	}

	return
}
