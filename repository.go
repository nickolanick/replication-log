package main

import (
//  "sync/atomic"
//  "sync"
)

type Repository struct {
  // total ordering
	messages []string
}

// TODO: implement RWLock
func (r *Repository) GetMessages() []string {
	return r.messages
}

// TODO: implement RWLock
func (r *Repository) AppendMessage(m string) {
	r.messages = append(r.messages, m)
}

var repository = Repository{[]string{}}
