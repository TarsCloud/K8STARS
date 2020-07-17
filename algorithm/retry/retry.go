package retry

import (
	"time"
)

// Option defines the strategy for retries
type Option interface {
	Init()
	CanRetry() bool
}

// OptionFactory is the retry option factory
type OptionFactory func() Option

// Func is a function for retries
type Func func(execFunc) error

type execFunc func() error

// New creates new retry strategy
func New(r OptionFactory) Func {
	return func(fun execFunc) error {
		var firstErr error
		opt := r()
		opt.Init()
		for {
			err := fun()
			if err == nil {
				return nil
			}
			if firstErr == nil {
				firstErr = err
			}
			if !opt.CanRetry() {
				return firstErr
			}
		}
	}
}

// MaxRetiesOpt returns a retry option factory
func MaxRetiesOpt(retryCount int) OptionFactory {
	return func() Option {
		return &maxRetiesOpt{
			retryCount: retryCount,
		}
	}

}

// MaxTimeoutOpt returns a retry option factory
func MaxTimeoutOpt(maxTimeout time.Duration, tryInterval time.Duration) OptionFactory {
	return func() Option {
		return &maxTimeoutOpt{
			maxTimeout:  maxTimeout,
			tryInterval: tryInterval,
		}
	}
}

type maxRetiesOpt struct {
	retryCount int
}

func (m *maxRetiesOpt) Init() {
}

func (m *maxRetiesOpt) CanRetry() bool {
	m.retryCount--
	if m.retryCount >= 0 {
		return true
	}
	return false
}

type maxTimeoutOpt struct {
	maxTimeout  time.Duration
	tryInterval time.Duration

	nextRetry time.Time
	deadline  time.Time
}

func (m *maxTimeoutOpt) Init() {
	now := time.Now()
	m.nextRetry = now.Add(m.tryInterval)
	m.deadline = now.Add(m.maxTimeout)
}

func (m *maxTimeoutOpt) CanRetry() bool {
	now := time.Now()
	m.nextRetry = now.Add(m.tryInterval)
	if now.After(m.deadline) {
		return false
	}
	if now.Before(m.nextRetry) {
		time.Sleep(m.nextRetry.Sub(now))
	}
	return true
}
