# gener8
simple golang go:generate code generator

# install
```
go get github.com/dbreedt/gener8
```

# supported parameters
```
skip_format: skip gofmt being run on the generated file
trace      : enables trace logging
in         : file to parse
out        : file to write the generated code to
pkg        : the value to replace $pkg with
kws        : csv list of values to replace $kwn tokens with
```

# usage
The best way to understand how to use it, is to look at the example in example folder.
The example demonstrates a simple redis functionality that can be generated and auto extended into a service.

This uses a non go file `redis.algo`
```go
/*
kw1: The name of the struct you want to extend
kw2: The sub-system name
kw3: The data type name
*/

package $pkg

import (
  "context"
  "encoding/json"
	"fmt"
  "log"

  "github.com/go-redis/redis/v8"
)

func (s *$kw1) redis$kw3Key(key int64) string {
	return fmt.Sprintf("$kw2_$kw3_%d", key)
}

func (s *$kw1) redisGet$kw3(ctx context.Context, key int64) *$kw3 {
  rKey := s.redis$kw3Key(key)
	data, err := s.redis.Get(ctx, rKey).Result()
	if err != nil {
    if err != redis.Nil {
      log.Println("Redis call failed for key", rKey)
    }
		return nil
	}

  item := &$kw3{}
  err = json.Unmarshal([]byte(data), item)
	if err != nil {
    log.Println("Redis::Unmarshal failed for key", rKey)
  }
	return item
}
```

When you run `go generate ./...` it will generate a file called `auto_generated_redis_sample.go` that stems from the `go:generate` directive in the `redis_smaple.go` file:
```go
//go:generate gener8 -in=redis.algo -out=auto_generated_redis_sample.go -pkg=example -kws=printerService,HP,Printer
```

The main gain from using this is that if we have several services like the printerService that we would like to have a common set of functionality/algorithms that work on different types, templates/generics would sure be handy here, we can write one algo file that contains the algorithm and use `go generate` like the c pre-compiler to generate all the "template" code for us. That way we have compile time type checking and no runtime surprises.

And it removes a lot of copy+pasta errors:
```go
func (s *printerService) GetPrinter(ctx context.Context, id int64) *Printer {
	item := s.redisGetPrinter(ctx, id)
	if item != nil {
		return item
	}
```