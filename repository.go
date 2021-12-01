package main

import (
  "sync"
  "container/heap"
  //"fmt"
)

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
// TODO: test - commit message to follower directly (out of range)
func (r *Repository) AppendMessage(wr_cons_msg WriteConsistencyMessage) {
  // first it should go to staging area
  // and then we try to pop and compare last element we have in current list with total order
  // only workaround is first element
	heap.Init(&r.staging_repo)
	heap.Push(&r.staging_repo, wr_cons_msg)

  fmt.Println(wr_cons_msg)
  for {
    if len(r.staging_repo) == 0 {
      break
    }
    wr_cons_msg := heap.Pop(&r.staging_repo).(WriteConsistencyMessage)
    if r.messages[len(r.messages)-1].TotalOrder + 1 == wr_cons_msg.TotalOrder {
	    r.messages = append(r.messages, wr_cons_msg)
    } else {
      heap.Push(&r.staging_repo, wr_cons_msg)
      break
    }
  }

  // fmt.Println(r.staging_repo)
  // TODO: try to pop and compare with last element in messages total order
}

var repository = Repository{[]WriteConsistencyMessage{{"none", 0, 0, &sync.WaitGroup{}}}, StagingRepository{}}
