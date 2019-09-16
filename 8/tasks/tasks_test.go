package tasks

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

// Counter
type counter struct {
	value int
	mx    *sync.RWMutex
}

// New counter struct (create mutext inside)
func newCounter(value int) *counter {
	return &counter{
		value: value,
		mx:    &sync.RWMutex{},
	}
}

// Concurrent-safe read counter
func (r *counter) read() int {
	r.mx.RLock()
	value := r.value
	r.mx.RUnlock()
	return value
}

// Concurrent-safe inrement counter
func (r *counter) add(n int) {
	r.mx.Lock()
	r.value += n
	r.mx.Unlock()
}

// Concurrent-safe set new value to counter
func (r *counter) set(v int) {
	r.mx.Lock()
	r.value += v
	r.mx.Unlock()
}

// Test return of run on empty slice of tasks
func TestNoTasks(t *testing.T) {
	count := Run(nil, 10, 10)
	if count != 0 {
		t.Errorf("run on empty slice must return 0 instread of %d", count)
	}
}

// Test return of run on one task without error
func TestOneTaskWithoutError(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return nil
		}),
	}
	count := Run(tasks, 10, 10)
	if count != 0 {
		t.Errorf("run on one task without error must return 0 instread of %d", count)
	}
}

// Test return of run on one task and it will end with error
func TestOneTaskWithError(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return errors.New("task failed")
		}),
	}
	count := Run(tasks, 10, 10)
	if count != 1 {
		t.Errorf("run on one task with error must return 1 instread of %d", count)
	}
}

// Test return of run on slice with 2 tasks without errors
func TestTwoTaskWithoutErrors(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return nil
		}),
	}
	count := Run(tasks, 10, 10)
	if count != 0 {
		t.Errorf("run on two tasks without errors must return 0 instread of %d", count)
	}
}

// Test return of run on slice with 2 tasks with errors
func TestTwoTaskWithErrors(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return errors.New("task failed")
		}),
		Task(func() error {
			return errors.New("task failed")
		}),
	}
	count := Run(tasks, 10, 10)
	if count != 2 {
		t.Errorf("run on two tasks with errors must return 2 instread of %d", count)
	}
}

// Test return of run on slice with 2 tasks: one with error, another wihout
func TestTwoMixedTasks(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return errors.New("task failed")
		}),
		Task(func() error {
			return nil
		}),
	}
	count := Run(tasks, 10, 10)
	if count != 1 {
		t.Errorf("run one task with error and one task without error must return 1 instread of %d", count)
	}
}

// Test return of run on slice of tasks (error and without) when number of concurency is 1
func TestRunTasksByOne(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return errors.New("task #1 failed")
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return errors.New("task #3 failed")
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return errors.New("task #5 failed")
		}),
	}

	count := Run(tasks, 1, 10)
	if count != 3 {
		t.Errorf("result must be 3 instread of %d", count)
	}
}

// Test return of run on slice of tasks (error and without) when number of concurency is 2
func TestRunTasksByTwo(t *testing.T) {
	tasks := []Task{
		Task(func() error {
			return errors.New("task #1 failed")
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return errors.New("task #3 failed")
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return errors.New("task #5 failed")
		}),
		Task(func() error {
			return errors.New("task #6 failed")
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return nil
		}),
		Task(func() error {
			return errors.New("task #9 failed")
		}),
		Task(func() error {
			return errors.New("task #10 failed")
		}),
	}

	count := Run(tasks, 2, 20)
	if count != 6 {
		t.Errorf("result must be 6 instread of %d", count)
	}
}

// Test that tasks run concurrently and utilize all available workers
// Test takes some time
func TestRunTasksConcurrently(t *testing.T) {

	// counter to tracking number of task running at the same time
	runningTasksCounter := newCounter(0)

	// counter that will tell that all tasks are done
	doneTasksCounter := newCounter(0)

	// quick task constructor
	newQuickTask := func(id int) Task {
		return Task(func() error {

			runningTasksCounter.add(1)

			time.Sleep(time.Duration(250) * time.Millisecond)

			runningTasksCounter.add(-1)
			doneTasksCounter.add(1)

			return nil
		})
	}

	// middle on duration task constructor
	newMiddleTask := func(id int) Task {
		return Task(func() error {

			runningTasksCounter.add(1)

			time.Sleep(time.Duration(500) * time.Millisecond)

			runningTasksCounter.add(-1)
			doneTasksCounter.add(1)

			return nil
		})
	}

	// long task constructor
	newLongTask := func(id int) Task {
		return Task(func() error {

			runningTasksCounter.add(1)

			time.Sleep(time.Duration(1) * time.Second)

			runningTasksCounter.add(-1)
			doneTasksCounter.add(1)

			return nil
		})
	}

	// work scenario
	tasks := []Task{
		newQuickTask(1),
		newLongTask(2),
		newMiddleTask(3),
		newMiddleTask(4),
		newQuickTask(5),
		newLongTask(6),
		newLongTask(7),
		newQuickTask(8),
		newMiddleTask(9),
		newLongTask(10),
		newQuickTask(11),
		newQuickTask(12),
		newLongTask(13),
		newMiddleTask(14),
		newQuickTask(15),
		newQuickTask(16),
		newMiddleTask(17),
		newQuickTask(18),
		newQuickTask(19),
		newLongTask(20),
	}

	tasksCount := len(tasks)

	// concurrency number
	n := 3

	go Run(tasks, n, -1)

	// peak (max) count of tasks running at a time
	peakRunningTasksCount := newCounter(0)

	// measure values of running tasks count
	var runningTasksCounts []int

	// measure ticker
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	// measure loop
	for range ticker.C {
		runningTasksCounts = append(runningTasksCounts, runningTasksCounter.read())

		// calc peam (max) measure
		if runningTasksCounter.read() > peakRunningTasksCount.read() {
			peakRunningTasksCount.set(runningTasksCounter.read())
		}

		// after all done stop measure loop (clean ticker)
		if doneTasksCounter.read() >= tasksCount {
			ticker.Stop()
			break
		}
	}

	// peak number of tasks running must not be greater than our concurrency number
	if peakRunningTasksCount.read() > n {
		t.Errorf("peak tasks count (%d) exceed %d", peakRunningTasksCount.read(), n)
	}

	// frequencies of measure values (peak and other all together);
	// peak measure must be more frequent, cause runner must fully utilize all available workers (=n)
	peakCountValueFrequency := 0
	otherCountValuesFrequency := 0

	for _, value := range runningTasksCounts {
		if peakRunningTasksCount.read() == value {
			peakCountValueFrequency++
		} else {
			otherCountValuesFrequency++
		}
	}

	if peakCountValueFrequency <= otherCountValuesFrequency {
		msg := `tasks count %d at a time happened less or equal than other count at a time.
Tasks count %d happened %d times
Other counts all together happened %d times.
Looks like runner not fully utilizes all available workers`
		t.Errorf(
			msg,
			peakRunningTasksCount.read(),
			peakRunningTasksCount.read(),
			peakCountValueFrequency,
			otherCountValuesFrequency,
		)
	}

}

// Test how works limit of fails
// Test takes some time
func TestRunTasksConcurrentlyWithLimit(t *testing.T) {

	// concurency number
	n := 4

	// limit of fails
	limit := 4

	// work time of one task (in ms)
	workTime := 250

	// successfull tasks will send own id after been worked,
	okCh := make(chan int, 100)

	// failed tasks will send own id after been worked,
	failCh := make(chan int, 100)

	// successfull task constructor
	newOkTask := func(id int) Task {
		return Task(func() error {
			time.Sleep(time.Duration(workTime) * time.Millisecond)
			okCh <- id
			return nil
		})
	}

	// fail task constructor
	newFailTask := func(id int) Task {
		return Task(func() error {
			time.Sleep(time.Duration(workTime) * time.Millisecond)
			err := fmt.Errorf("task #%d failed", id)
			failCh <- id
			return err
		})
	}

	// Prepare work scenario
	// Failure limit is 4, so key is 5th failure task
	// To time this task done runner could possible already had get 3 more tasks in running state (#14, #15, #16)
	tasks := []Task{

		newFailTask(1),
		newOkTask(2),
		newOkTask(3),
		newFailTask(4),

		newOkTask(5),
		newFailTask(6),
		newOkTask(7),
		newOkTask(8),

		newFailTask(9),
		newFailTask(10), // <-- 5th failure task
		newOkTask(11),   // it is possible could be executed after 5th failure
		newOkTask(12),   // it is possible could be executed after 5th failure

		newOkTask(13),   // it is possible could be executed after 5th failure
		newFailTask(14), // it is possible could be executed after 5th failure
		newOkTask(15),   // it is possible could be executed after 5th failure
		newOkTask(16),   // it is possible could be executed after 5th failure

		newOkTask(17),   // unlikly could be executed after 5th failure
		newFailTask(18), // ...
		newOkTask(19),
		newOkTask(20),

		newOkTask(21),
		newFailTask(22),
		newOkTask(23),
		newOkTask(24),
	}

	// Run tasks
	go Run(tasks, n, limit)

	// total count of all tasks
	tasksCount := len(tasks)

	// calculate total worktime (+500 ms just for sure)
	totalWorkTime := (workTime*tasksCount)/n + 500

	time.Sleep(time.Duration(totalWorkTime) * time.Millisecond)

	// read ids from fail channel
	failChLen := len(failCh)
	runFailTaskIds := make([]int, failChLen)
	for i := 0; i < failChLen; i++ {
		v := <-failCh
		runFailTaskIds[i] = v
	}

	expectedRunFailTaskIds1 := []int{1, 4, 6, 9, 10}
	expectedRunFailTaskIds2 := []int{1, 4, 6, 9, 10, 14}

	// need to sort so we can use equal
	sort.Ints(runFailTaskIds)

	// limit is 4, to time 5th fail task (#10) done, runner could possible already had get 3 more tasks in running state (#14,#15,#16)
	if !reflect.DeepEqual(runFailTaskIds, expectedRunFailTaskIds1) &&
		!reflect.DeepEqual(runFailTaskIds, expectedRunFailTaskIds2) {
		t.Errorf("could expect %v or %v, instread of %v", expectedRunFailTaskIds1, expectedRunFailTaskIds2, runFailTaskIds)
	}

	// read ids from ok channel
	okChLen := len(okCh)
	runOkTaskIds := make([]int, okChLen)
	for i := 0; i < okChLen; i++ {
		v := <-okCh
		runOkTaskIds[i] = v
	}

	// need to sort so we can use equal
	sort.Ints(runOkTaskIds)

	expectedRunOkTasksIds1 := []int{2, 3, 5, 7, 8, 11}
	expectedRunOkTasksIds2 := []int{2, 3, 5, 7, 8, 11, 12}
	expectedRunOkTasksIds3 := []int{2, 3, 5, 7, 8, 11, 12, 13}
	expectedRunOkTasksIds4 := []int{2, 3, 5, 7, 8, 11, 12, 13, 15}
	expectedRunOkTasksIds5 := []int{2, 3, 5, 7, 8, 11, 12, 13, 15, 16}

	// limit is 4, to time 5th fail task (#10) done, runner could possible already had get 3 more tasks in running state (#14,#15,#16)
	if !reflect.DeepEqual(runOkTaskIds, expectedRunOkTasksIds1) &&
		!reflect.DeepEqual(runOkTaskIds, expectedRunOkTasksIds2) &&
		!reflect.DeepEqual(runOkTaskIds, expectedRunOkTasksIds3) &&
		!reflect.DeepEqual(runOkTaskIds, expectedRunOkTasksIds4) &&
		!reflect.DeepEqual(runOkTaskIds, expectedRunOkTasksIds5) {
		t.Errorf("could expect %v or %v or %v or %v or %v, instread of %v", expectedRunOkTasksIds1, expectedRunOkTasksIds2,
			expectedRunOkTasksIds3, expectedRunOkTasksIds4, expectedRunOkTasksIds5, runOkTaskIds)
	}
}
