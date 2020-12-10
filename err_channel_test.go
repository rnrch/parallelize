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
