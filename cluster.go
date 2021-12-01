package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
)

type Cluster struct {
	nodes []Node
}

var cluster = Cluster{[]Node{}}

func (cluster *Cluster) init(followers []string) {
	for _, follower := range followers {
		msg_queue := make(chan *WriteConsistencyMessage, 200000)
		node := Node{follower, msg_queue, "healthy"}
		go node.commitMessage()
		cluster.nodes = append(cluster.nodes, node)
	}
}

func (cluster *Cluster) commitMessages(wcmsg *WriteConsistencyMessage) { // for chan in chans send message
	for _, node := range cluster.nodes {
		node.msg_queue <- wcmsg
	}
}

// TODO: cluster.status (reuse in quorum)

type Node struct {
	addr      string
	msg_queue chan *WriteConsistencyMessage
	// TODO: rewrite with enum
	health string
}

func (node *Node) healthCheck() {
  // run ping
  // implement health API -> customer facing, should proxy request to primary if secondary
  // implement ping API -> non-customer facing
  // during cluster init -> same as commitMessage()
  // health -> should be enum (red, yellow green)
  // uint, 0 (red) 3 (green) other -> yellow red[0 yellow 3]green
}

// TODO: change name
func (node *Node) commitMessage() {

	for wcmessage := range node.msg_queue {
		// send request to commit
		postBody, _ := json.Marshal(wcmessage)

		for {
      // TODO: use health information to increase retry sleep delay in for loop
      // kill/pause one of container
			responseBody := bytes.NewBuffer(postBody)
			_, err := http.Post(node.addr+"/commit", "application/json", responseBody)
			// TODO: read response body status and retry
			if err == nil {
				break
			}
		}

		// negative wait group throws panic
		// we can catch it inside another function scope
		func(cond *sync.WaitGroup) {
			// for some reason defer recover() does not work
			defer func() { recover() }()
			cond.Done()
		}(wcmessage.WriteCond)
	}
}
