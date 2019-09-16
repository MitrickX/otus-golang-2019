package tasks

import (
	"sync"
)

// Task type
type Task func() error

// Concurrent-safe counter with limit - limiter
// If limit exceed current counter isExceed will tell about it
type limiter struct {
	count int // current limiter
	limit int // limit
	mx    *sync.RWMutex
}

// New limiter struct (create mutext inside)
func newLimiter(limit int) *limiter {
	return &limiter{
		mx:    &sync.RWMutex{},
		limit: limit,
	}
}

// Concurrent-safe inrement limiter
func (r *limiter) add(n int) {
	r.mx.Lock()
	r.count += n
	r.mx.Unlock()
}

// Concurrent-safe read current count of limiter
func (r *limiter) read() int {
	r.mx.RLock()
	count := r.count
	r.mx.RUnlock()
	return count
}

// Is limiter exceed his limit?
func (r *limiter) isExceed() bool {
	r.mx.RLock()
	exceed := r.count > r.limit
	r.mx.RUnlock()
	return exceed
}

// Helper worker
// Read task from tasks channel and run it
// On error result of task increment limiter
// If tasks channel closed (range is done), decrement wait group
func runWorker(tasks <-chan Task, limiter *limiter, wg *sync.WaitGroup) {
	for task := range tasks {
		if !limiter.isExceed() {
			err := task()
			if err != nil {
				limiter.add(1)
			}
		}
	}
	wg.Done()
}

// Run tasks with concurrency number and fails limit
//
// - n:
// 		Number of tasks that could be run concurrently (at the same time).
//     	Runner fully utilize all available goroutines-workers (<= n) at a time
//   n == 0 means no tasks will running
//   n < 0  means all tasks will running, same as n == len(tasks)
//
// - limit:
//     	Number of max errors that allowed be happen before runner will stop
//     	After runner have stop some tasks could be possible already in running state.
//	   	So they will work until stops by itself, and Run will wait them
//   limit <= 0 means limit will not taking into account (there will not stopping by limit), same as limit = len(tasks)
//
// Returns number of fails happended
func Run(tasks []Task, n int, limit int) int {

	// boundary case
	if n == 0 {
		return 0
	}

	// total number of tasks
	tasksCount := len(tasks)

	// boundary case
	if tasksCount == 0 {
		return 0
	}

	// n < 0: number of concurrency is number of tasks
	if n < 0 {
		n = tasksCount
	}

	// limit <= 0: no limit (or in another words, limit is number of tasks)
	if limit <= 0 {
		limit = tasksCount
	}

	// number of concurrency must not be bigger number of tasks
	if n > tasksCount {
		n = tasksCount
	}

	// channel with tasks
	tasksCh := make(chan Task)

	// limiter
	limiter := newLimiter(limit)

	// run n workes
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go runWorker(tasksCh, limiter, wg)
	}

	// feed workers while there are tasks and limit is not exceed
	for i := 0; i < tasksCount && !limiter.isExceed(); i++ {
		tasksCh <- tasks[i]
	}

	close(tasksCh)

	// wait until all currently running tasks stop
	wg.Wait()

	// number of error happened
	return limiter.read()
}
