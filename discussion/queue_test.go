package discussion

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
	"time"
)

/*
Cases:
- passing in a nil pointer
- passing a topic with nil ("") name
- normal case, name doesn't already exist
- attempting to add a new topic of the same name
*/
func TestQueue_Add(t *testing.T) {
	q := NewQueue()

	tests := []struct {
		name    string
		in      *Topic
		wantLen int
		wantErr bool
	}{
		{name: "nil-topic", in: nil, wantLen: 0, wantErr: true},
		{name: "empty-name", in: &Topic{Name: ""}, wantLen: 0, wantErr: true},
		{name: "normal", in: &Topic{Name: "test"}, wantLen: 1, wantErr: false},
		{name: "dup-name", in: &Topic{Name: "test"}, wantLen: 1, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := q.Add(test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if test.wantLen != q.Len() {
				t.Errorf("len != wantLen (len = %d, wantLen = %d)", q.Len(), test.wantLen)
			}
		})
	}
}

/*
Cases:
- topic with that name not found
- normal case
*/
func TestQueue_Attach(t *testing.T) {
	q := NewQueue()
	_ = q.Add(&Topic{
		Name:        "test2",
		Description: "blah",
		Sources:     nil,
		Modified:    time.Now(),
	})

	tests := []struct {
		name    string
		inName  string
		inStr   string
		wantErr bool
	}{
		{name: "name-not-found", inName: "test1", inStr: "google.com", wantErr: true},
		{name: "normal", inName: "test2", inStr: "google.com", wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := q.Attach(test.inName, test.inStr)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if !test.wantErr {
				if q.Q[0].Sources[0] != test.inStr {
					t.Errorf("sources didn't get updated")
				}
			}
		})
	}
}

/*
Cases:
- No topics, topic not found
- One topic
- 5 Topics: back to front, front to front, mid to front, mid to front
*/
func TestQueue_Bump(t *testing.T) {
	q := NewQueue()

	// zero tops; not found
	err := q.Bump("some-name")
	if err == nil {
		t.Errorf("no error thrown")
	}

	// one top
	_ = q.Add(&Topic{Name: "t1"})
	err = q.Bump("t1")
	if q.Q[0].Name != "t1" {
		t.Errorf("bump moved the only topic back")
	}

	_ = q.Add(&Topic{Name: "t2"})
	_ = q.Add(&Topic{Name: "t3"})
	_ = q.Add(&Topic{Name: "t4"})
	_ = q.Add(&Topic{Name: "t5"})

	tests := []struct {
		name      string
		in        string
		wantOrder []string
	}{
		{name: "back-to-front", in: "t5", wantOrder: []string{"t5", "t1", "t2", "t3", "t4"}},
		{name: "front", in: "t5", wantOrder: []string{"t5", "t1", "t2", "t3", "t4"}},
		{name: "mid", in: "t3", wantOrder: []string{"t3", "t5", "t1", "t2", "t4"}},
		{name: "one-to-front", in: "t1", wantOrder: []string{"t1", "t3", "t5", "t2", "t4"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = q.Bump(test.in)
			for i := range test.wantOrder {
				if q.Q[i].Name != test.wantOrder[i] {
					t.Errorf("incorrect order (i = %d), (want = %s, got = %s)", i, test.wantOrder[i], q.Q[i].Name)
				}
			}
		})
	}
}

/*
Cases:
- Name not found
- Source not found
- Normal case(s)
- Nil sources slice
*/
func TestQueue_Detach(t *testing.T) {
	q := NewQueue()
	_ = q.Add(&Topic{
		Name:        "testTopic",
		Description: "for testing",
		Sources:     []string{"google.com", "twitter.com", "amazon.com"},
		Modified:    time.Now(),
	})

	tests := []struct {
		name      string
		topicName string
		drop      int
		wantSrc   []string
		wantErr   bool
	}{
		{name: "topic-not-found", topicName: "none", drop: 0, wantSrc: nil, wantErr: true},
		{name: "source-not-found", topicName: "testTopic", drop: 4, wantSrc: nil, wantErr: true},
		{name: "normal1", topicName: "testTopic", drop: 1, wantSrc: []string{"google.com", "amazon.com"}, wantErr: false},
		{name: "normal2", topicName: "testTopic", drop: 0, wantSrc: []string{"amazon.com"}, wantErr: false},
		{name: "normal3", topicName: "testTopic", drop: 0, wantSrc: []string{}, wantErr: false},
		{name: "nil-source", topicName: "testTopic", drop: 0, wantSrc: nil, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := q.Detach(test.topicName, test.drop)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			// don't run more tests if we want an error
			if !test.wantErr {
				if !cmp.Equal(q.Q[0].Sources, test.wantSrc) {
					t.Errorf("sources slices don't agree")
				}
			}
		})
	}
}

/*
Cases:
- Nil topic slice
- Not found name
- Normal case
*/
func TestQueue_Next(t *testing.T) {
	q := NewQueue()

	_, err := q.Next()
	if err == nil {
		t.Errorf("expected an error, got a nil")
	}

	_ = q.Add(&Topic{Name: "t1"})

	top, err := q.Next()
	if err != nil {
		t.Errorf("expected a nil, got an error")
	}

	if top == nil {
		t.Errorf("expected a topic got a nil")
	}

	if top.Name != "t1" {
		t.Errorf("topic has wrong name")
	}
}

/*
Cases:
- Not found
- Remove the last item
- Remove the first item
- Remove a middle item
- Normal case
*/
func TestQueue_Remove(t *testing.T) {
	q := NewQueue()
	_ = q.Add(&Topic{Name: "test1"})
	_ = q.Add(&Topic{Name: "test2"})
	_ = q.Add(&Topic{Name: "test3"})
	_ = q.Add(&Topic{Name: "test4"})
	_ = q.Add(&Topic{Name: "test5"})

	tests := []struct {
		name    string
		in      string
		wantLen int
		wantErr bool
	}{
		{name: "not-found", in: "not-there", wantLen: 5, wantErr: true},
		{name: "remove-last", in: "test5", wantLen: 4, wantErr: false},
		{name: "remove-first", in: "test1", wantLen: 3, wantErr: false},
		{name: "remove-middle", in: "test3", wantLen: 2, wantErr: false},
		{name: "remove1", in: "test2", wantLen: 1, wantErr: false},
		{name: "remove2", in: "test4", wantLen: 0, wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := q.Remove(test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if !test.wantErr {
				if q.Len() != test.wantLen {
					t.Errorf("len != wantLen (len = %d, wantLen = %d)", q.Len(), test.wantLen)
				}

				for _, topic := range q.Q {
					if topic.Name == test.in {
						t.Errorf("a name was not removed (%s)", test.in)
					}
				}
			}
		})
	}
}

/*
Cases:
- non existent
- normal cases
*/
func TestQueue_Skip(t *testing.T) {
	q := NewQueue()

	// nil q
	err := q.Skip("non-extant")
	if err == nil {
		t.Errorf("skip added a Topic")
	}

	// one item
	_ = q.Add(&Topic{Name: "t1"})
	err = q.Skip("t1")
	if q.Len() != 1 {
		t.Errorf("skip added a topic")
	}

	if q.Q[0].Name != "t1" {
		t.Errorf("skip moved t1 when it shouldn't have")
	}

	_ = q.Add(&Topic{Name: "t2"})
	_ = q.Add(&Topic{Name: "t3"})
	_ = q.Add(&Topic{Name: "t4"})
	_ = q.Add(&Topic{Name: "t5"})

	tests := []struct {
		name      string
		in        string
		wantOrder []string
	}{
		{name: "last", in: "t5", wantOrder: []string{"t1", "t2", "t3", "t4", "t5"}},
		{name: "first", in: "t1", wantOrder: []string{"t2", "t3", "t4", "t5", "t1"}},
		{name: "middle", in: "t4", wantOrder: []string{"t2", "t3", "t5", "t1", "t4"}},
		{name: "second", in: "t3", wantOrder: []string{"t2", "t5", "t1", "t4", "t3"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = q.Skip(test.in)
			for i := range test.wantOrder {
				if q.Q[i].Name != test.wantOrder[i] {
					t.Errorf("names don't match at index %d (want = %s, got = %s)", i, test.wantOrder[i], q.Q[i].Name)
				}
			}
		})
	}
}

/*
Test Cases:
- normal
*/
func TestQueue_SaveLoad(t *testing.T) {
	dir := "./test"

	// check if the dir exists, create if it doesn't
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			t.FailNow()
		}
	}

	testQ := NewQueue()
	err := testQ.Add(&Topic{
		Name:        "testName",
		Description: "testDescription",
		Sources: []string{
			"https://www.google.com/",
			"https://www.wikipedia.org/",
		},
		Modified:  time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC),
		Created:   time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC),
		CreatedBy: "Tester",
	})
	if err != nil {
		t.FailNow()
	}

	err = testQ.Save("./test")
	if err != nil {
		t.Errorf("Error on save: %v", err)
	}

	newQ := NewQueue()
	err = newQ.Load("./test")
	if err != nil {
		t.Errorf("Error on load: %v", err)
	}

	if !cmp.Equal(testQ, newQ) {
		t.Errorf("testQ != newQ:\n%s", cmp.Diff(testQ, newQ))
	}

	// delete dir and all files
	err = os.RemoveAll(dir)
	if err != nil {
		t.FailNow()
	}
}
