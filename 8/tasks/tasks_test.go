package tasks

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

// Test return of function that run on empty slice of tasks
func TestNoTasks(t *testing.T) {
	count := Run(nil, 10, 10)
	if count != 0 {
		t.Errorf("run on empty slice must return 0 instread of %d", count)
	}
}

// Test return of function that run on one task without error
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

// Test return of function that run on one task and it will end with error
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

// Test return of function that run on slice with 2 tasks without errors
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

// Test return of function that run on slice with 2 tasks with errors
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

// Test return of function that run on slice with 2 tasks: one with error, another wihout
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

// Test return of function that run on slice of tasks (error and without) when number of concurency is 1
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

// Test return of function that run on slice of tasks (error and without) when number of concurency is 2
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

// Test that 2 tasks run concurrently (not test result of function, just concurrency)
// Test on slice of 20 tasks
// Test without limit (all tasks work without errors)
func TestRunTasksConcurrently2(t *testing.T) {
	testRunTasksConcurrently(t, 2)
}

// Test that 4 tasks run concurrently (not test result of function, just concurrency)
// Test on slice of 20 tasks
// Test without limit (all tasks work without errors)
func TestRunTasksConcurrently4(t *testing.T) {
	testRunTasksConcurrently(t, 4)
}

// Test that 8 tasks run concurrently (not test result of function, just concurrency)
// Test on slice of 20 tasks
// Test without limit (all tasks work without errors)
func TestRunTasksConcurrently8(t *testing.T) {
	testRunTasksConcurrently(t, 8)
}

// Test that all tasks run concurrently (not test result of function, just concurrency)
// Test on slice of 20 tasks
// Test without limit (all tasks work without errors)
func TestRunTasksConcurrentlyAll(t *testing.T) {
	testRunTasksConcurrently(t, -1)
}

// Test how works limit of fails (not test result of function, not test of concurrency, just limitation)
// Test on slice of 16 tasks
func TestRun16TasksConcurrently4AndWithLimit4(t *testing.T) {

	// concurency number
	n := 4

	// limit of fails
	limit := 4

	// successfull tasks will send own id after been worked,
	okCh := make(chan int, 100)

	// failed tasks will send own id after been worked,
	failCh := make(chan int, 100)

	// successfull task constructor
	newOkTask := func(id int) Task {
		return Task(func() error {
			time.Sleep(time.Duration(1) * time.Second)
			okCh <- id
			return nil
		})
	}

	// fail task constructor
	newFailTask := func(id int) Task {
		return Task(func() error {
			time.Sleep(time.Duration(1) * time.Second)
			err := fmt.Errorf("task #%d failed", id)
			failCh <- id
			return err
		})
	}

	// we have 4 bunch of tasks, in each bunch we have 4 of concurently running tasks
	// After 5th failed task runner already should not run new tasks (tasks from another banch),
	// but another tasks in current banch could be worked or could be not (cause we already run goroutine for these tasks)
	tasks := []Task{

		// 1st bunch of concurently running tasks
		newFailTask(1),
		newOkTask(2),
		newOkTask(3),
		newFailTask(4),

		// 2nd bunch of concurently running tasks
		newOkTask(5),
		newFailTask(6),
		newOkTask(7),
		newOkTask(8),

		// 3d bunch of concurently running tasks
		newFailTask(9),
		newFailTask(10),
		newOkTask(11), // task could be worked or could be not
		newOkTask(12), // task could be worked or could be not

		// 4th bunch of concurently running tasks (all tasks for sure must not be worked)
		newOkTask(13),
		newFailTask(14),
		newOkTask(15),
		newOkTask(16),
	}

	// Run tasks
	go Run(tasks, n, limit)

	// bunch tasks takes 1 s (cause of concurently)
	// 4 bunches takes 4 s (bunckes run is sequential) + 500ms for for sure
	time.Sleep(time.Duration(4500) * time.Millisecond)

	// read ids from fail channel
	failChLen := len(failCh)
	runFailTaskIds := make([]int, failChLen)
	for i := 0; i < failChLen; i++ {
		v := <-failCh
		runFailTaskIds[i] = v
	}

	expectedRunFailTaskIds := []int{1, 4, 6, 9, 10}

	// limit is 4, after 5th fail task (#10) runner must not run another bunch, next (6th) fail task (#14) is never run
	//   => len == 5
	if !isEqualsIntSlicesAsSets(runFailTaskIds, expectedRunFailTaskIds) {
		t.Errorf("expected %v instread of %v", expectedRunFailTaskIds, runFailTaskIds)
	}

	// read ids from ok channel
	okChLen := len(okCh)
	runOkTaskIds := make([]int, okChLen)
	for i := 0; i < okChLen; i++ {
		v := <-okCh
		runOkTaskIds[i] = v
	}

	expectedRunOkTasks1Ids := []int{2, 3, 5, 7, 8}
	expectedRunOkTasks2Ids := []int{2, 3, 5, 7, 8, 11}
	expectedRunOkTasks3Ids := []int{2, 3, 5, 7, 8, 11, 12}

	// after fail task #10 we could has worked tasks #11, #12, but not from another bunch
	ok := isEqualsIntSlicesAsSets(runOkTaskIds, expectedRunOkTasks1Ids) ||
		isEqualsIntSlicesAsSets(runOkTaskIds, expectedRunOkTasks2Ids) ||
		isEqualsIntSlicesAsSets(runOkTaskIds, expectedRunOkTasks3Ids)
	if !ok {
		t.Errorf("expected %v or %v or %v, instread of %v", expectedRunOkTasks1Ids, expectedRunOkTasks2Ids,
			expectedRunOkTasks3Ids, runOkTaskIds)
	}
}

// Private helper for test concurrently
// Test that n of tasks run concurrently
// - tasks is slice of tasks
// - n is number of concurrency, -1 run all tasks concurrently
// Notice that this test take some time
func testRunTasksConcurrently(t *testing.T, n int) {

	// structure for tracking start times of tasks
	type taskStartTime struct {
		id    int
		start int64 // UnixNano results
	}

	// tasks will send start times
	startTimeCh := make(chan taskStartTime, 100)

	// task constructor
	newTask := func(id int) Task {
		return Task(func() error {
			time.Sleep(time.Duration(1) * time.Second)
			startTimeCh <- taskStartTime{id, time.Now().UnixNano()}
			return nil
		})
	}

	// prepare tasks list
	taskCount := 20
	tasks := make([]Task, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = newTask(i + 1)
	}

	if n < 0 {
		n = taskCount
	}

	// Run tasks
	go Run(tasks, n, -1)

	// calculate number of bunches
	bunchNum := taskCount / n
	if taskCount%n > 0 {
		bunchNum++
	}

	// Tasks in each bunch takes 1 s, so also a bunch (cause of concurently)
	// Sleep whole time for all bunches + 500 ms just for sure
	time.Sleep(time.Duration(bunchNum*1000+500) * time.Millisecond)

	// read start times from channel
	startTimeChLen := len(startTimeCh)
	startTimes := make([]taskStartTime, startTimeChLen)
	for i := 0; i < startTimeChLen; i++ {
		startTimes[i] = <-startTimeCh
	}

	// sort by id of tasks, so we can correct testing each bunch of timings
	sort.Slice(startTimes, func(i, j int) bool {
		return startTimes[i].id < startTimes[j].id
	})

	// slice startTimes by n
	for offset := 0; offset < taskCount; offset += n {

		// id of bunch (slice)
		bunchId := (offset / n) + 1

		// calc low:high of slice
		low := offset
		high := offset + n
		if high >= taskCount {
			high = taskCount
		}

		// start task times of bunch
		times := startTimes[low:high]

		// min, max of the times
		var minTime, maxTime int64
		for _, item := range times {
			if minTime == 0 || item.start < minTime {
				minTime = item.start
			}
			if maxTime == 0 || item.start > maxTime {
				maxTime = item.start
			}
		}

		// all tasks in bunch run concurrently and must started in time range of 1 second (cause each task in bunch concurrent and work 1 second)
		diff := maxTime - minTime
		if diff > int64(time.Second) {
			diffS := float64(diff) / float64(time.Second)
			t.Errorf("bunch #%d of tasks must start concurrently in time range of 1 second, instread of %0.3f s", bunchId, diffS)
		}
	}
}

// Check if 2 int slices are equals without mattering values' order
func isEqualsIntSlicesAsSets(a []int, b []int) bool {
	a = cloneIntSlice(a)
	b = cloneIntSlice(b)
	sort.Ints(a)
	sort.Ints(b)
	return reflect.DeepEqual(a, b)
}

// Clone int slice
func cloneIntSlice(a []int) []int {
	b := make([]int, len(a))
	copy(a, b)
	return b
}
