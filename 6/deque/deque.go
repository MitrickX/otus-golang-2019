package deque

// Linked list node struct
type dequeIntNode struct {
	value int
	next  *dequeIntNode
	prev  *dequeIntNode
}

// Linked list node struct
type DequeInt struct {
	first *dequeIntNode
	last  *dequeIntNode
	n     int
}

// If deque empty
func (deque *DequeInt) IsEmpty() bool {
	return deque.n == 0
}

// How much items in deque
func (deque *DequeInt) Size() int {
	return deque.n
}

// Add value to the front
func (deque *DequeInt) AddFirst(value int) {

	// Linked list node
	node := &dequeIntNode{
		value: value,
	}

	if deque.IsEmpty() {
		// init case: first and last must point to this new node
		deque.first = node
		deque.last = node
	} else {
		// insert new node to the front of list and move first pointer to this new node
		node.next = deque.first
		deque.first.prev = node
		deque.first = node
	}

	// increment Size of linded list
	deque.n++
}

// Add the item to the back
func (deque *DequeInt) AddLast(value int) {
	// Linked list node
	node := &dequeIntNode{
		value: value,
	}

	if deque.IsEmpty() {
		// init case: first and last must point to this new node
		deque.first = node
		deque.last = node
	} else {
		// insert new node to the back of list and move last pointer to this new node
		node.prev = deque.last
		deque.last.next = node
		deque.last = node
	}

	// increment Size of linded list
	deque.n++
}

// Remove and return value from the front
// If removing is OK - return TRUE and value
// Otherwise coudn't remove cause deque is empty - return FALSE and empty value
func (deque *DequeInt) RemoveFirst() (int, bool) {

	// boundary case - empty
	if deque.IsEmpty() {
		return 0, false
	}

	// boundary case: first and last is the same
	if deque.n == 1 {
		value := deque.first.value
		deque.first = nil
		deque.last = nil
		deque.n--
		return value, true
	}

	// now we has at least 2 nodes in list

	// prev first node
	prevFirst := deque.first

	// new first node (move to next node)
	deque.first = deque.first.next

	// unlink new first node with prev node (prevFirst)
	deque.first.prev = nil

	// unlink prevFirst with new first
	prevFirst.next = nil

	// deincrement Size
	deque.n--

	return prevFirst.value, true
}

// Remove and return the value from the back
// If removing is OK - return TRUE and value
// Otherwise coudn't remove cause deque is empty - return FALSE and empty value
func (deque *DequeInt) RemoveLast() (int, bool) {

	// boundary case - empty
	if deque.IsEmpty() {
		return 0, false
	}

	// boundary case: first and last is the same
	if deque.n == 1 {
		value := deque.first.value
		deque.first = nil
		deque.last = nil
		deque.n--
		return value, true
	}

	// now we has at least 2 nodes in list

	// prev last node
	prevLast := deque.last

	// new last node (move to prev node)
	deque.last = deque.last.prev

	// unlink new last node with next node (prevLast)
	deque.last.next = nil

	// unlink prevLast with new last
	prevLast.prev = nil

	// deincrement Size
	deque.n--

	return prevLast.value, true
}

// Read value from the front
// If reading is OK - return TRUE and value
// Otherwise coudn't read cause deque is empty - return FALSE and empty value
func (deque *DequeInt) ReadFirst() (int, bool) {
	// boundary case - empty
	if deque.IsEmpty() {
		return 0, false
	}
	return deque.first.value, true
}

// Read value from the back
// If reading is OK - return TRUE and value
// Otherwise coudn't read cause deque is empty - return FALSE and empty value
func (deque *DequeInt) ReadLast() (int, bool) {
	// boundary case - empty
	if deque.IsEmpty() {
		return 0, false
	}
	return deque.last.value, true
}

// Read all values in deque
func (deque *DequeInt) ReadAll() []int {
	// boundary case - empty
	if deque.IsEmpty() {
		return nil
	}

	n := deque.n

	// allocate result slice with proper size
	result := make([]int, n)

	// on each step collect current node value and move to next node
	node := deque.first
	for i := 0; i < n; i++ {
		result[i] = node.value
		node = node.next
	}

	return result
}

func NewDequeInt() *DequeInt {
	return &DequeInt{}
}
