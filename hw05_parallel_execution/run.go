package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		m = 1
	}

	var currentTaskIdx int32 = -1
	var errCounter int32

	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	once := sync.Once{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-doneCh:
					return
				default:
				}
				idx := atomic.AddInt32(&currentTaskIdx, 1)
				if int(idx) >= len(tasks) {
					break
				}
				err := tasks[idx]()
				if err != nil {
					numErrs := atomic.AddInt32(&errCounter, 1)

					if int(numErrs) >= m {
						once.Do(func() {
							close(doneCh)
						})
						break
					}
				}
			}
		}()
	}

	wg.Wait()

	if int(errCounter) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
