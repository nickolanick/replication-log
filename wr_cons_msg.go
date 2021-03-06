package main

import (
	"encoding/json"
	"sync"
)

type WriteConsistencyMessage struct {
	Message string `json:"message"`
	// write consistency
	WriteConsistency int `json:"write_consistency"`
	TotalOrder       int `json:"total_order"`
	WriteCond        *sync.WaitGroup
}

type WriteConsistencyMessageJSON struct {
	Message string `json:"message"`
	// write consistency
	WriteConsistency int `json:"write_consistency"`
	TotalOrder       int `json:"total_order,omitempty"`
}

func (wr_cons_msg *WriteConsistencyMessage) UnmarshalJSON(data []byte) error {
	var res WriteConsistencyMessageJSON

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	wr_cons_msg.Message = res.Message
	wr_cons_msg.TotalOrder = res.TotalOrder

	// TODO: change to something else or check if default provided
	// user can potentially pass with 0 to not wait for commit
	wr_cons_msg.WriteConsistency = res.WriteConsistency

	if wr_cons_msg.WriteConsistency == 0 {
		wr_cons_msg.WriteConsistency = len(config.nodes)
	}

	wr_cons_msg.WriteCond = &sync.WaitGroup{}
	wr_cons_msg.WriteCond.Add(wr_cons_msg.WriteConsistency)

	return nil
}
