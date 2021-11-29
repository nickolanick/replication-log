package main

import (
	//  "sync/atomic"
	"encoding/json"
	"sync"
)

type WriteConsistencyMessage struct {
	Message string
	// write consistency
	WriteConsistency int
	WriteCond        *sync.WaitGroup
}

type WriteConsistencyMessageJSON struct {
	Message string `json:"message"`
	// write consistency
	WriteConsistency int `json:"write_consistency"`
}

func (wr_cons_msg *WriteConsistencyMessage) UnmarshalJSON(data []byte) error {
	var res WriteConsistencyMessageJSON

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	wr_cons_msg.Message = res.Message

	// TODO: change to something else or check if default provided
	// user can potentially pass with 0 to not wait for commit
	wr_cons_msg.WriteConsistency = res.WriteConsistency

	if wr_cons_msg.WriteConsistency == 0 {
		wr_cons_msg.WriteConsistency = len(config.nodes)
	}

	var write_cond sync.WaitGroup

	write_cond.Add(wr_cons_msg.WriteConsistency)
	wr_cons_msg.WriteCond = &write_cond

	return nil
}
