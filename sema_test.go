package main

import (
	"math"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func sleep1s() error {
	time.Sleep(time.Second)
	return nil
}

func TestSemaphore(t *testing.T) {
	sema := NewSemaphore(4) // set max currency to 4
	now := time.Now()
	for i := 0; i < 4; i++ {
		sema.Go(sleep1s)
	}
	err := sema.Wait()
	assert.Equal(t, err, nil)
	sec := math.Round(time.Since(now).Seconds())
	assert.Equal(t, int(sec), 1)

	sema = NewSemaphore(2) // set max currency to 2
	now = time.Now()
	for i := 0; i < 4; i++ {
		sema.Go(sleep1s)
	}
	err = sema.Wait()
	assert.Equal(t, err, nil)
	sec = math.Round(time.Since(now).Seconds())
	assert.Equal(t, int(sec), 2)
}
