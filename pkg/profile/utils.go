package profile

import (
	"context"
	"log"
	"time"
)

// Интерфейс оборачиваемой ф-ции для Retry
type Effoctor func(context.Context) error

// Реализация паттерна Retry
// Если ресурс недоступен, то пробуем подергать его еще
func Retry(effoctor Effoctor, retries int, delay time.Duration) Effoctor {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			log.Println("...Attempt ", r)
			err := effoctor(ctx)
			if err == nil || r > retries {
				return err
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				log.Println("canceling")
				return ctx.Err()
			}

		}
	}
}
