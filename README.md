# Interrupt

A simple Go library for turning OS signals into context cancellations


```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/relvacode/interrupt"
)


func main() {
    // By default the context cancels on os.Interrupt
    ctx := interrupt.Context(context.Background())

    select {
        case <-ctx.Done():
            fmt.Println(ctx.Err())
        case <-time.After(time.Second * 30):
            fmt.Println("not signalled")
    }
}
```