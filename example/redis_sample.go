//go:generate gener8 -in=redis.algo -out=auto_generated_redis_sample.go -pkg=example -kws=printerService,HP,Printer
package example

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Printer struct {
	Name          string
	Make          string
	Model         string
	DriverVersion string
}

type printerService struct {
	redis *redis.Client
}

func NewPrinterService() *printerService {
	return &printerService{
		redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (s *printerService) GetPrinter(ctx context.Context, id int64) *Printer {
	item := s.redisGetPrinter(ctx, id)
	if item != nil {
		return item
	}

	// TOOD :: add db call here
	// TODO :: add redisSetPrinter call here

	return nil
}
