package parallelize

import "context"

// ErrorChannel supports non-blocking send and receive operation to capture error.
// A maximum of one error is kept in the channel and the rest of the errors sent
// are ignored, unless the existing error is received and the channel becomes empty
// again.
type ErrorChannel struct {
	ch chan error
}

// NewErrorChannel returns a new ErrorChannel.
func NewErrorChannel() *ErrorChannel {
	return &ErrorChannel{
		ch: make(chan error, 1),
	}
}

// SendError sends an error without blocking the sender.
func (e *ErrorChannel) SendError(err error) {
	select {
	case e.ch <- err:
	default:
	}
}

// SendErrorWithCancel sends an error without blocking the sender and calls
// cancel function.
func (e *ErrorChannel) SendErrorWithCancel(err error, cancel context.CancelFunc) {
	e.SendError(err)
	cancel()
}

// ReceiveError receives an error from channel without blocking on the receiver.
func (e *ErrorChannel) ReceiveError() error {
	select {
	case err := <-e.ch:
		return err
	default:
		return nil
	}
}
