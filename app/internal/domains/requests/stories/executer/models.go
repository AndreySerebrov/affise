package executer

import "time"

type Config struct {
	RequestsInParallel int64
	MaxUrlNum          int64
	RequestTimeout     time.Duration
}

type Response struct {
	Url      string
	Response string
}

type ResponseList struct {
	List []Response
}
