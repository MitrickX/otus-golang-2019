package tasks

import (
	"sync"
)

// Task type
type Task func() error

// Counter
type counter struct {
	count int // current counter
	limit int // limit
	mx    *sync.RWMutex
}

// New counter struct (create mutext inside)
func newCounter(limit int) *counter {
	return &counter{
		mx:    &sync.RWMutex{},
		limit: limit,
	}
}

// Concurrent-safe read counter
func (r *counter) Count() int {
	r.mx.RLock()
	count := r.count
	r.mx.RUnlock()
	return count
}

// Concurrent-safe inrement counter
func (r *counter) incCount(n int) {
	r.mx.Lock()
	r.count += n
	r.mx.Unlock()
}

// Is counter exceed his limit
func (r *counter) isExceed() bool {
	r.mx.RLock()
	exceed := r.count > r.limit
	r.mx.RUnlock()
	return exceed
}

// Helper to run bunch of tasks - each task run concurrently
// - cntr counter
// If counter exceed stop run tasks
// Function is blocking, using wait group so function works until all tasks stops
func run(tasks []Task, cntr *counter) {

	wg := &sync.WaitGroup{}

	for _, task := range tasks {

		// stop run tasks
		if cntr.isExceed() {
			break
		}

		// increment wait group
		wg.Add(1)
		go func(task Task) {
			if !cntr.isExceed() {
				err := task()
				if err != nil {
					cntr.incCount(1)
				}
			}
			wg.Done()
		}(task)

	}

	wg.Wait()
}

// Run tasks with concurrency number and fails limit
//
// - n: number of tasks that could be run concurrently
//   n == 0 means no tasks will running
//   n < 0  means all tasks will running, same as n == len(tasks)
//
// - limit: number of max errors that could be happen before stop running
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

	// goroutine safe counter with limit checker
	// counter will count number of errors
	couner := newCounter(limit)

	// go though all tasks, slice by n size and run concurrently
	for offset := 0; offset < tasksCount; offset += n {

		// low bound of slice
		low := offset

		// high bound of slice, must not be bigger number of tasks
		high := offset + n
		if high >= tasksCount {
			high = tasksCount
		}

		// slice (bunch) of tasks to run concurrently
		bunch := tasks[low:high]

		// run concurrently
		run(bunch, couner)

		// stop working
		if couner.isExceed() {
			break
		}
	}

	// Total count of errors happend
	return couner.Count()
}
