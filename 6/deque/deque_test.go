package deque

import (
	"reflect"
	"testing"
)

// Test IsEmpty correctness on new deque and after adding and removing
func TestIsEmpty(t *testing.T) {
	deque := NewDequeInt()
	AssertTrue(t, deque.IsEmpty(), "New deque must be empty")

	deque.AddFirst(1)
	AssertFalse(t, deque.IsEmpty(), "Deque with after addFirst must be not empty")

	deque.AddLast(2)
	AssertFalse(t, deque.IsEmpty(), "Deque with after addFirst and addLast must be not empty")

	deque.RemoveLast()
	AssertFalse(t, deque.IsEmpty(), "Deque with after addFirst, addLast, RemoveLast must be not empty")

	deque.RemoveFirst()
	AssertTrue(t, deque.IsEmpty(), "Deque with after addFirst, addLast, RemoveLast, RemoveFirst must be empty")
}

// Test Size correctness on new deque and after adding and removing
func TestSize(t *testing.T) {
	deque := NewDequeInt()
	AssertIntEquals(t, 0, deque.Size(), "New deque must has 0 size")

	deque.AddFirst(10)
	AssertIntEquals(t, 1, deque.Size(), "Deque with after one addFirst must has size 1")

	deque.AddLast(22)
	AssertIntEquals(t, 2, deque.Size(), "Deque with after one addFirst and one addLast must has size 2")

	deque.RemoveLast()
	AssertIntEquals(t, 1, deque.Size(), "Deque with after addFirst, addLast, RemoveLast must has size 1")

	deque.RemoveFirst()
	AssertIntEquals(t, 0, deque.Size(), "Deque with after addFirst, addLast, RemoveLast, RemoveFirst must has size 0")
}

// Test ReadAll correctness on new deque and after adding and removing
func TestReadAll(t *testing.T) {
	deque := NewDequeInt()

	values := deque.ReadAll()
	AssertLen(t, 0, values, "ReadAll on new deque must return empty slice (len=0)")

	deque.AddLast(1)
	deque.AddLast(2)
	deque.AddFirst(3)
	deque.AddFirst(4)
	deque.AddLast(5)

	values = deque.ReadAll()

	AssertSliceEquals(t, []int{4, 3, 1, 2, 5}, values, "ReadAll after 5 calls of add methods must return [4, 3, 1, 2, 5]")

	// removing calls
	deque.RemoveLast()
	deque.RemoveFirst()

	values = deque.ReadAll()

	AssertSliceEquals(t, []int{3, 1, 2}, values, "ReadAll after 5 calls of adding and 2 calls of removing methods must return [3, 1, 2]")

	// drain rest items
	deque.RemoveFirst()
	deque.RemoveLast()
	deque.RemoveLast()

	values = deque.ReadAll()

	AssertLen(t, 0, values, "ReadAll after dain all values must return empty slice (len=0)")
}

// Test read first value on new deque and after adding and removing
func TestReadFirst(t *testing.T) {
	deque := NewDequeInt()

	var ok bool
	var value int

	_, ok = deque.ReadFirst()

	AssertFalse(t, ok, "Read first from new (empty) deque must not be ok")

	deque.AddLast(1)

	value, ok = deque.ReadFirst()

	AssertTrue(t, ok, "Read first after 1 call of adding must be ok")

	AssertIntEquals(t, 1, value, "Read first after AddLast(1) must be return 1")

	deque.AddFirst(2)

	value, ok = deque.ReadFirst()

	AssertTrue(t, ok, "Read first after 2 calls of adding must be ok")

	AssertIntEquals(t, 2, value, "Read first after call addFirst(2) must be return 2")

	deque.AddLast(3)

	value, ok = deque.ReadFirst()

	AssertTrue(t, ok, "Read first after 3 calls of adding must be ok")

	AssertIntEquals(t, 2, value, "Read first after calls - addFirst(2) and AddLast(3) must be return 2")

	deque.RemoveLast()

	value, ok = deque.ReadFirst()

	AssertTrue(t, ok, "Read first after 3 calls of adding and 1 call of removing must be ok")

	AssertIntEquals(t, 2, value, "Read first after RemoveLast() after which deque is [2 1] must be return 2")

	deque.RemoveFirst()

	value, ok = deque.ReadFirst()

	AssertTrue(t, ok, "Read first after 3 calls of adding and 2 calls of removing must be ok")

	AssertIntEquals(t, 1, value, "Read first after RemoveLast() after whith deque is [1] must be return 1")
}

// Test read last value on new deque and after adding and removing
func TestReadLast(t *testing.T) {
	deque := NewDequeInt()

	var ok bool
	var value int

	_, ok = deque.ReadLast()

	AssertFalse(t, ok, "Read last from new (empty) deque must not be ok")

	deque.AddFirst(1)

	value, ok = deque.ReadLast()

	AssertTrue(t, ok, "Read last after 1 call of adding must be ok")

	AssertIntEquals(t, 1, value, "Read last after AddFirst(1) must be return 1")

	deque.AddLast(2)

	value, ok = deque.ReadLast()

	AssertTrue(t, ok, "Read last after 2 calls of adding must be ok")

	AssertIntEquals(t, 2, value, "Read first after call AddLast(2) must be return 2")

	deque.AddFirst(3)

	value, ok = deque.ReadLast()

	AssertTrue(t, ok, "Read last after 3 calls of adding must be ok")

	AssertIntEquals(t, 2, value, "Read last after calls - AddLast(2) and AddFirst(3) must be return 2")

	deque.RemoveFirst()

	value, ok = deque.ReadLast()

	AssertTrue(t, ok, "Read last after 3 calls of adding and 1 call of removing must be ok")

	AssertIntEquals(t, 2, value, "Read last after RemoveFirst() after which deque is [1,2] must be return 2")

	deque.RemoveLast()

	value, ok = deque.ReadLast()

	AssertTrue(t, ok, "Read last after 3 calls of adding and 2 calls of removing must be ok")

	AssertIntEquals(t, 1, value, "Read first after RemoveLast() after whith deque is [1] must be return 1")

}

// Test remove first on new deque and after some adding
func TestRemoveFirst(t *testing.T) {
	deque := NewDequeInt()

	var ok bool
	var value int

	_, ok = deque.RemoveLast()

	AssertFalse(t, ok, "Remove from new (empty) deque must not be ok")

	deque.AddLast(1)
	deque.AddFirst(2)
	deque.AddLast(3)
	deque.AddLast(4)
	deque.AddFirst(5)

	value, ok = deque.RemoveFirst()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 5, value, "Remove first from [5, 2, 1, 3, 4] must be 5")

	value, ok = deque.RemoveFirst()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 2, value, "Remove first from [2, 1, 3, 4] must be 2")

	value, ok = deque.RemoveFirst()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 1, value, "Remove first from [1, 3, 4] must be 1")

	value, ok = deque.RemoveFirst()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 3, value, "Remove first from [3, 4] must be 3")

	value, ok = deque.RemoveFirst()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 4, value, "Remove first from [4] must be 4")

	_, ok = deque.RemoveFirst()

	AssertFalse(t, ok, "Remove from empty deque must be not ok")
}

// Test remove last on new deque and after some adding
func TestRemoveLast(t *testing.T) {
	deque := NewDequeInt()

	var ok bool
	var value int

	_, ok = deque.RemoveLast()

	AssertFalse(t, ok, "Remove from new (empty) deque must not be ok")

	deque.AddLast(1)
	deque.AddFirst(2)
	deque.AddLast(3)
	deque.AddLast(4)
	deque.AddFirst(5)

	value, ok = deque.RemoveLast()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 4, value, "Remove first from [5, 2, 1, 3, 4] must be 4")

	value, ok = deque.RemoveLast()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 3, value, "Remove first from [5, 2, 1, 3] must be 3")

	value, ok = deque.RemoveLast()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 1, value, "Remove first from [5, 2, 1] must be 1")

	value, ok = deque.RemoveLast()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 2, value, "Remove first from [5, 2] must be 2")

	value, ok = deque.RemoveLast()

	AssertTrue(t, ok, "Remove from not empty deque must be ok")
	AssertIntEquals(t, 5, value, "Remove first from [5] must be 5")

	_, ok = deque.RemoveLast()

	AssertFalse(t, ok, "Remove from empty deque must be not ok")
}

func AssertTrue(t *testing.T, expr bool, message string) {
	if !expr {
		t.Errorf("Expected to be TRUE: %s\n", message)
	}
}

func AssertFalse(t *testing.T, expr bool, message string) {
	if expr {
		t.Errorf("Expected to be FALSE: %s\n", message)
	}
}

func AssertIntEquals(t *testing.T, expected int, tested int, message string) {
	if expected != tested {
		t.Errorf("Expected %d not equals tested %d : %s\n", expected, tested, message)
	}
}

func AssertLen(t *testing.T, length int, tested []int, message string) {
	if len(tested) != length {
		t.Errorf("Expected length of slice is %d not %d: %s\n", length, len(tested), message)
	}
}

func AssertSliceEquals(t *testing.T, expected []int, tested []int, message string) {
	if !reflect.DeepEqual(expected, tested) {
		t.Errorf("Expected %v not equals tested %v: %s\n", expected, tested, message)
	}
}
