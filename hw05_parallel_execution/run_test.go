package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestLessZeroErrors(t *testing.T) {
	t.Run("test m < 0 value", func(t *testing.T) {
		const taskCount = 55
		var executed int32

		tasks := make([]Task, 0, taskCount)

		for i := 0; i < taskCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&executed, 1)
				return nil
			})
		}

		err := Run(tasks, 10, -55)
		require.NoError(t, err)
		require.Equal(t, int32(taskCount), atomic.LoadInt32(&executed))
	})
	t.Run("test m == 0 value", func(t *testing.T) {
		const taskCount = 55
		var executed int32

		tasks := make([]Task, 0, taskCount)

		for i := 0; i < taskCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&executed, 1)
				return nil
			})
		}

		err := Run(tasks, 10, 0)
		require.NoError(t, err)
		require.Equal(t, int32(taskCount), atomic.LoadInt32(&executed))
	})
}

func TestRun_ErrorsLimitExceeded(t *testing.T) {
	const tasksCount = 1000
	var executed int32

	tasks := make([]Task, 0, tasksCount)
	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&executed, 1)
			return errors.New("task error")
		})
	}

	n, m := 10, 3
	err := Run(tasks, n, m)
	require.Error(t, err)
	require.Equal(t, ErrErrorsLimitExceeded, err)
	require.LessOrEqual(t, atomic.LoadInt32(&executed), int32(n+m))
}

func TestConcurrencyWithoutSleep(t *testing.T) {
	const numTasks = 50
	const NumNeededGorut = 50
	var currentGorutines int32

	tasks := make([]Task, numTasks)

	waitingCh := make(chan struct{})

	for i := 0; i < numTasks; i++ {
		tasks[i] = func() error {
			atomic.AddInt32(&currentGorutines, 1)

			<-waitingCh
			atomic.AddInt32(&currentGorutines, -1)
			return nil
		}
	}

	var err error
	doneRunCh := make(chan struct{})
	go func() {
		err = Run(tasks, NumNeededGorut, 1)
		close(doneRunCh)
	}()

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&currentGorutines) == int32(NumNeededGorut)
	}, 1*time.Second, 10*time.Millisecond)

	close(waitingCh)
	<-doneRunCh
	require.NoError(t, err)
	require.Equal(t, atomic.LoadInt32(&currentGorutines), int32(0))
}
