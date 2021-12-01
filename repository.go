package main

import (
  "sync"
  "container/heap"
  "fmt"
)

// TODO: implement some way to test it
type Repository struct {
  // total ordering
	messages []WriteConsistencyMessage
  staging_repo StagingRepository
}

// TODO: implement RWLock
func (r *Repository) GetMessages() []string {
  var msgs []string
  // this is little hack to not overcomplicte code in append message with conditionals
  for _, msg := range r.messages[1:] {
        msgs = append(msgs, msg.Message)
  }
	return msgs
}

// TODO: implement RWLock
func (r *Repository) AppendMessage(wr_cons_msg WriteConsistencyMessage) {
  // first it should go to staging area
  // and then we try to pop and compare last element we have in current list with total order
  // only workaround is first element
	heap.Init(&r.staging_repo)
	heap.Push(&r.staging_repo, wr_cons_msg)
  for {
    // get last element in heap
    // if repository last element total order is exactly one less than arrived
    // than append
    // otherwise break
  }

  fmt.Println(r.staging_repo)
  // TODO: try to pop and compare with last element in messages total order
	r.messages = append(r.messages, wr_cons_msg)
}

var repository = Repository{[]WriteConsistencyMessage{{"none", 0, -1, &sync.WaitGroup{}}}, StagingRepository{}}
