package discussion

import (
	"encoding/json"
	"fmt"
	"github.com/ericebersohl/gobottas/discord"
	"io/ioutil"
	"log"
	"time"
)

// Slice for simplicity, no need to make it a heap-based PQ
type Queue struct {
	Q        []*Topic  `json:"q"`        // hide the internal list from the user
	Modified time.Time `json:"modified"` // time last modified
}

// Create a new Queue, initializes the underlying slice and updates Modified
func NewQueue() *Queue {
	q := Queue{
		Q:        make([]*Topic, 0),
		Modified: time.Now(),
	}
	return &q
}

// Get the number of topics in the queue
func (q *Queue) Len() int {
	return len(q.Q)
}

// Return all topics in the Queue
func (q *Queue) List() []*Topic {
	return q.Q
}

// Return the first topic in the queue.  Does not remove the topic from the queue
func (q *Queue) Next() (*Topic, error) {
	if len(q.Q) > 0 {
		return q.Q[0], nil
	}
	return nil, discord.NewError("Empty Queue", "Cannot call next when the queue is empty.")
}

// Add a Topic to the queue
func (q *Queue) Add(t *Topic) error {

	// check for nil
	if t == nil {
		return discord.NewError("Nil Topic", "Cannot add a nil topic to the queue.")
	}

	// check for invalid name
	if t.Name == "" {
		return discord.NewError("Empty Topic Name", "Cannot add a topic with no name.")
	}

	// check for name that already exists
	for _, topic := range q.Q {
		if topic.Name == t.Name {
			return discord.NewError("Duplicate Topic", "A topic with that name already exists.")
		}
	}

	// append to the list
	q.Q = append(q.Q, t)

	// update modified
	q.Modified = time.Now()

	return nil
}

// Removes a Topic of the specified name from the Queue, does nothing if the name is not found
func (q *Queue) Remove(s string) error {

	found := false

	// find the name
	for i, t := range q.Q {

		// remove the name
		if t.Name == s {
			q.Q = append(q.Q[:i], q.Q[i+1:]...)

			// update modified
			q.Modified = time.Now()

			// set found
			found = true
		}
	}

	q.Modified = time.Now()

	if found == false {
		return discord.NewError("Topic Not Found", "Could not find a topic with that name.")
	}
	return nil
}

// Moves the Topic of the specified name to the front of the Queue
func (q *Queue) Bump(s string) error {

	found := false

	// find the Topic
	for i := range q.Q {
		if q.Q[i].Name == s {

			// pull out the topic
			t := q.Q[i]

			// rebuild the slice and prepend the topic
			q.Q = append(q.Q[:i], q.Q[i+1:]...)
			q.Q = append([]*Topic{t}, q.Q...)

			q.Q[0].Modified = time.Now()

			found = true
		}
	}

	q.Modified = time.Now()

	if found == false {
		return discord.NewError("Topic Not Found", "Could not find a topic with that name.")
	}
	return nil
}

// moves the specified Topic to the end of the Queue
func (q *Queue) Skip(s string) error {

	found := false

	// find the topic
	for i := range q.Q {
		if q.Q[i].Name == s {
			tmp := q.Q[i]
			tmp.Modified = time.Now()
			q.Q = append(q.Q[:i], q.Q[i+1:]...)
			q.Q = append(q.Q, tmp)

			found = true
		}
	}

	q.Modified = time.Now()

	if found == false {
		return discord.NewError("Topic Not Found", "Could not find a topic with that name.")
	}
	return nil
}

// attach a string to the list of sources
func (q *Queue) Attach(n, s string) error {
	found := false

	for _, t := range q.Q {
		if t.Name == n {
			t.Sources = append(t.Sources, s)
			t.Modified = time.Now()
			found = true
		}
	}

	q.Modified = time.Now()

	if found == false {
		return discord.NewError("Topic Not Found", "Could not find a topic with that name.")
	}
	return nil
}

// remove a source (by index) from the specified topic
func (q *Queue) Detach(n string, i int) error {
	found := false

	for _, t := range q.Q {
		if t.Name == n {
			if len(t.Sources) <= i || i < 0 {
				return discord.NewError("Index Out of Range", "You specified a number that is out of the range of sources.")
			}
			fmt.Println(t.Sources)
			t.Sources = append(t.Sources[:i], t.Sources[i+1:]...)
			t.Modified = time.Now()
			found = true
		}
	}

	q.Modified = time.Now()

	if found == false {
		return discord.NewError("Topic Not Found", "Could not find a topic with that name.")
	}

	return nil
}

// Save the state of the queue to JSON
func (q *Queue) Save(path string) error {
	// get []byte
	data, err := json.Marshal(q)
	if err != nil {
		log.Printf("Save error: %v", err)
		return err
	}

	// write to file
	err = ioutil.WriteFile(fmt.Sprintf("%s/queue.json", path), data, 0644)
	if err != nil {
		log.Printf("WriteFile: %v", err)
		return err
	}

	return nil
}

// load data into queue from JSON
func (q *Queue) Load(path string) error {

	// get data
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/queue.json", path))
	if err != nil {
		log.Printf("Load: %v", err)
		return err
	}

	// unmarshal data
	err = json.Unmarshal(data, q)
	if err != nil {
		log.Printf("Load Unmarshal error: %v", err)
		return err
	}

	return nil
}
