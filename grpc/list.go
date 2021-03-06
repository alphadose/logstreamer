package grpc

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/alphadose/logstreamer/types"
)

// global memory pool for storing and leasing node objects
var (
	nodePool = sync.Pool{New: func() any { return unsafe.Pointer(new(node)) }}
	nodeGet  = nodePool.Get
	nodePut  = nodePool.Put
	store    = NewList() // main data storage for GRPC server
)

// List is a lock-free linked list
// theory -> https://www.cs.rochester.edu/u/scott/papers/1996_PODC_queues.pdf
// pseudocode -> https://www.cs.rochester.edu/research/synchronization/pseudocode/queues.html
// much faster and efficient than container/list + sync.Mutex
// Used as a thread-safe in-memory storage for payload objects by the GRPC server
type List struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// NewList returns a new list
func NewList() List {
	n := nodeGet().(unsafe.Pointer)
	(*node)(n).next, (*node)(n).value = nil, nil
	return List{head: n, tail: n}
}

// a single node in the linked list
type node struct {
	next  unsafe.Pointer
	value *types.Payload
}

// Enqueue inserts a value into the list
func (l *List) Enqueue(value *types.Payload) {
	var (
		n          = nodeGet().(unsafe.Pointer)
		tail, next unsafe.Pointer
	)
	(*node)(n).next, (*node)(n).value = nil, value
	for {
		tail = atomic.LoadPointer(&l.tail)
		next = atomic.LoadPointer(&(*node)(tail).next)
		if tail == atomic.LoadPointer(&l.tail) { // are tail and next consistent?
			if next == nil {
				if atomic.CompareAndSwapPointer(&(*node)(tail).next, next, n) {
					atomic.CompareAndSwapPointer(&l.tail, tail, n) // Enqueue is done.  try to swing tail to the inserted node
					return
				}
			} else { // tail was not pointing to the last node
				// try to swing Tail to the next node
				atomic.CompareAndSwapPointer(&l.tail, tail, next)
			}
		}
	}
}

// Dequeue removes and returns the value at the head of the queue to the memory pool
// It returns nil if the list is empty
func (l *List) Dequeue() (value *types.Payload) {
	var head, tail, next unsafe.Pointer
	for {
		head = atomic.LoadPointer(&l.head)
		tail = atomic.LoadPointer(&l.tail)
		next = atomic.LoadPointer(&(*node)(head).next)
		if head == atomic.LoadPointer(&l.head) { // are head, tail, and next consistent?
			if head == tail { // is list empty or tail falling behind?
				if next == nil { // is list empty?
					return nil
				}
				// tail is falling behind.  try to advance it
				atomic.CompareAndSwapPointer(&l.tail, tail, next)
			} else {
				// read value before CAS_node otherwise another dequeue might free the next node
				value = (*node)(next).value
				if atomic.CompareAndSwapPointer(&l.head, head, next) {
					(*node)(head).next, (*node)(head).value = nil, nil
					nodePut(head)
					return // Dequeue is done.  return
				}
			}
		}
	}
}
