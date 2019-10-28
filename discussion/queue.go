package discussion

import (
	"errors"
	"fmt"
	"time"
)

// Slice for simplicity, no need to make it a heap-based PQ
type Queue struct {
	q        []*Topic // hide the internal list from the user
	Modified time.Time
}

// Create a new Queue, initializes the underlying slice and updates Modified
func NewQueue() *Queue {
	q := Queue{
		q:        make([]*Topic, 0),
		Modified: time.Now(),
	}
	return &q
}

// Get the number of topics in the queue
func (q *Queue) Len() int {
	return len(q.q)
}

// Return the first topic in the queue.  Does not remove the topic from the queue
func (q *Queue) Next() (*Topic, error) {
	if len(q.q) > 0 {
		return q.q[0], nil
	}
	return nil, errors.New("the queue is empty")
}

// Add a Topic to the queue
func (q *Queue) Add(t *Topic) error {

	// check for nil
	if t == nil {
		return errors.New("cannot Add nil topic")
	}

	// check for invalid name
	if t.Name == "" {
		return errors.New("cannot Add a topic with no name")
	}

	// check for name that already exists
	for _, topic := range q.q {
		if topic.Name == t.Name {
			return errors.New("a topic with that name already exists")
		}
	}

	// append to the list
	q.q = append(q.q, t)

	// update modified
	q.Modified = time.Now()

	return nil
}

// Removes a Topic of the specified name from the Queue, does nothing if the name is not found
func (q *Queue) Remove(s string) {

	// find the name
	for i, t := range q.q {

		// remove the name
		if t.Name == s {
			q.q = append(q.q[:i], q.q[i+1:]...)

			// update modified
			q.Modified = time.Now()
		}
	}
}

// Moves the Topic of the specified name to the front of the Queue
func (q *Queue) Bump(s string) {

	// find the Topic
	for i := range q.q {
		if q.q[i].Name == s {

			// pull out the topic
			t := q.q[i]

			// rebuild the slice and prepend the topic
			q.q = append(q.q[:i], q.q[i+1:]...)
			q.q = append([]*Topic{t}, q.q...)
		}
	}
}

// moves the specified Topic to the end of the Queue
func (q *Queue) Skip(s string) {
	// find the topic
	for i := range q.q {
		if q.q[i].Name == s {
			tmp := q.q[i]
			q.q = append(q.q[:i], q.q[i+1:]...)
			q.q = append(q.q, tmp)
		}
	}
}

// attach a string to the list of sources
func (q *Queue) Attach(n, s string) error {
	found := false

	for _, t := range q.q {
		if t.Name == n {
			t.Sources = append(t.Sources, s)
			found = true
		}
	}

	if found == false {
		return errors.New("the specified topic does not exist")
	}
	return nil
}

// remove a source (by index) from the specified topic
func (q *Queue) Detach(n string, i int) error {
	found := false

	for _, t := range q.q {
		if t.Name == n {
			if len(t.Sources) <= i || i < 0 {
				return errors.New("index out of range")
			}
			fmt.Println(t.Sources)
			t.Sources = append(t.Sources[:i], t.Sources[i+1:]...)
			found = true
		}
	}

	if found == false {
		return errors.New("the specified topic does not exist")
	}

	return nil
}
