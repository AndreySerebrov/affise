package loclimiter

type Limiter struct {
	ch chan bool
}

func New(limit int64) *Limiter {
	return &Limiter{
		ch: make(chan bool, limit),
	}
}

// Взять задачу
func (l *Limiter) Take() chan bool {
	return l.ch
}

// Задача выполнена
func (l *Limiter) Free() {
	<-l.ch
}
