package main

import (
//  "sync/atomic"
  "sync"
)

type WriteConsistencyMessage struct {
	Message          string `json:"message"`
	// write consistency
	WriteConsistency int  `json:"write_consistency"`
  WriteCond *sync.WaitGroup `json:"write_cond,omitempty"`
}
