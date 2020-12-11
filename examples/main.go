// Copyright 2020 rnrch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"

	"github.com/rnrch/parallelize"
	"github.com/rnrch/rlog"
)

func addTen(ctx context.Context, num int) (int, error) {
	rlog.Info("add 10 for number", "number", num)
	return num + 10, nil
}

func general() {
	rlog.Info("Start func", "case", "general")
	ctx, cancel := context.WithCancel(context.Background())
	errCh := parallelize.NewErrorChannel()
	raw := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := make([]int, len(raw))
	parallelize.Until(ctx, len(raw), func(index int) {
		res, err := addTen(ctx, raw[index])
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
}

func errOnEven(ctx context.Context, num int) error {
	if num%2 == 0 {
		return fmt.Errorf("got even number %d", num)
	}
	rlog.Info("Got odd number", "number", num)
	return nil
}

func stopOnError() {
	rlog.Info("Start func", "case", "stopOnError")
	ctx, cancel := context.WithCancel(context.Background())
	errCh := parallelize.NewErrorChannel()
	raw := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 2, 27, 29, 31, 33, 35, 37, 39}
	parallelize.Until(ctx, len(raw), func(index int) {
		err := errOnEven(ctx, raw[index])
		if err != nil {
			errCh.SendErrorWithCancel(err, cancel)
			return
		}
	}, parallelize.WithParallelism(2))
	if err := errCh.ReceiveError(); err != nil {
		rlog.Error(err, "Running stop on error case", "func", "errOnEven", "raw", raw)
	}
}

func continueOnError() {
	rlog.Info("Start func", "case", "continueOnError")
	ctx := context.Background()
	raw := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 2, 27, 29, 31, 33, 35, 37, 39}
	errCh := parallelize.NewErrorChannel()
	parallelize.Until(ctx, len(raw), func(index int) {
		err := errOnEven(ctx, raw[index])
		if err != nil {
			errCh.SendError(err)
			return
		}
	}, parallelize.WithParallelism(2))
	if err := errCh.ReceiveError(); err != nil {
		rlog.Error(err, "Running continue on error case", "func", "errOnEven", "raw", raw)
	}
}

func main() {
	general()
	stopOnError()
	continueOnError()
}
