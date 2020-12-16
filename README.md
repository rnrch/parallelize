# parallelize

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rnrch/parallelize)](https://pkg.go.dev/github.com/rnrch/parallelize)
[![Go Report Card](https://goreportcard.com/badge/github.com/rnrch/parallelize)](https://goreportcard.com/report/github.com/rnrch/parallelize)
![Github Actions](https://github.com/rnrch/parallelize/workflows/CI/badge.svg)

A minimal package to parallelize work and return error, with respect to k8s.io/client-go/util/workqueue

## Usage

e.g.

```go
    ctx, cancel := context.WithCancel(context.Background())
    errCh := parallelize.NewErrorChannel()
    raw := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    result := make([]int, len(raw))
    parallelize.Until(ctx, len(raw), func(index int) {
        res, err := addTen(ctx, raw[index]) // addTen returns (10 + input number)
        if err != nil {
            errCh.SendErrorWithCancel(err, cancel)
            return
        }
        result[index] = res
    })
    if err := errCh.ReceiveError(); err != nil {
        rlog.Error(err, "Running general case", "func", "addTen", "raw", raw)
    }
    rlog.Info("General case result", "func", "addTen", "raw", raw, "result", result)
```

see more in [examples](examples/main.go)
