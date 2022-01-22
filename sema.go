package main

import (
	"context"
	"sync"
)

type Semaphore struct {
	c       chan struct{}
	wg      sync.WaitGroup
	cancel  func()
	errOnce sync.Once
	err     error
}

func NewSemaphore(maxSize int) *Semaphore {
	return &Semaphore{
		c: make(chan struct{}, maxSize),
	}
}

func NewSemaphoreWithContext(ctx context.Context, maxSize int) (*Semaphore, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Semaphore{
		c:      make(chan struct{}, maxSize),
		cancel: cancel,
	}, ctx
}

func (s *Semaphore) Go(f func() error) {
	s.wg.Add(1)
	s.c <- struct{}{}
	go func() {
		defer func() {
			<-s.c
			s.wg.Done()
		}()
		if err := f(); err != nil {
			s.errOnce.Do(func() {
				s.err = err
				if s.cancel != nil {
					s.cancel()
				}
			})
		}
	}()
}

func (s *Semaphore) Wait() error {
	s.wg.Wait()
	if s.cancel != nil {
		s.cancel()
	}
	return s.err
}
