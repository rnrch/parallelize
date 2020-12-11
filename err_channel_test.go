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

package parallelize

import (
	"context"
	"errors"
	"testing"
)

func TestErrorChannel(t *testing.T) {
	ch := NewErrorChannel()

	if err := ch.ReceiveError(); err != nil {
		t.Errorf("expect nil from err channel, but got %v", err)
	}

	err := errors.New("unknown error")
	ch.SendError(err)
	if e := ch.ReceiveError(); e != err {
		t.Errorf("expect %v from err channel, but got %v", err, e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch.SendErrorWithCancel(err, cancel)
	if e := ch.ReceiveError(); e != err {
		t.Errorf("expect %v from err channel, but got %v", err, e)
	}

	if ctxErr := ctx.Err(); ctxErr != context.Canceled {
		t.Errorf("expect context canceled, but got %v", ctxErr)
	}
}
