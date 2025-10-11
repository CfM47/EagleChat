package datastructures

import (
	"container/list"
	"sync"
)

// Queue is a generic queue implementation based on container/list.
type Queue[T any] struct {
	container *list.List
	mutex     sync.Mutex
}

// NewQueue creates and returns a new Queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{container: list.New()}
}

// Enqueue adds a value to the back of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.container.PushBack(value)
}

// Dequeue removes and returns the value from the front of the queue.
// It returns the zero value of the type and false if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.container.Len() == 0 {
		var zero T
		return zero, false
	}
	element := q.container.Front()
	q.container.Remove(element)
	return element.Value.(T), true
}

// Peek returns the value at the front of the queue without removing it.
// It returns the zero value of the type and false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.container.Len() == 0 {
		var zero T
		return zero, false
	}
	return q.container.Front().Value.(T), true
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.container.Len() == 0
}

// Size returns the number of elements in the queue.
func (q *Queue[T]) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.container.Len()
}
