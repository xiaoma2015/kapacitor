package diagnostic

import (
	"bytes"
	"sync"
	"time"
)

type Data struct {
	Time    time.Time
	Message string
	Level   string
	Context []Field
	Fields  []Field
}

func (d Data) WriteTo(buf *bytes.Buffer) (int, error) {
	return 0, nil
}

type node struct {
	data *Data
	next *node
}

type Queue struct {
	head *node
	tail *node

	length int
	mu     sync.Mutex
}

func (q *Queue) Len() int {
	return q.length
}

func (q *Queue) Enqueue(d *Data) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.length < 0 {
		panic("queue length should never be less than 0")
	}

	q.length++

	n := &node{
		data: d,
	}

	if q.tail == nil {
		q.tail = n
		q.head = n

		return
	}

	q.tail.next = n
	q.tail = n

	return
}

func (q *Queue) Dequeue() *Data {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.length <= 0 {
		return nil
	}

	q.length--

	d := q.head.data
	// TODO: is this right??
	if q.length == 0 {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.next
	}

	return d
}