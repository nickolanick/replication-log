package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Cluster struct {
	nodes []Node
}

var health_status = map[int]string{
	0: "red",
	1: "yellow",
	2: "green",
}

var cluster = Cluster{[]Node{}}

func (cluster *Cluster) init(followers []string) {
	for _, follower := range followers {
		msg_queue := make(chan *WriteConsistencyMessage, 200000)
		health := Health{2, true}
		node := Node{follower, msg_queue, &health}

		if config.role == "leader" {
			go node.healthCheck()
		}
		go node.commitMessage()
		cluster.nodes = append(cluster.nodes, node)
	}
}

func (cluster *Cluster) commitMessages(wcmsg *WriteConsistencyMessage) { // for chan in chans send message
	for _, node := range cluster.nodes {
		node.msg_queue <- wcmsg
	}
}

func (cluster *Cluster) status() string {
	//node.health.status
	result := ""
	for _, node := range cluster.nodes {
		result += fmt.Sprintf("\nNode: %s, Status: %s", node.addr, health_status[node.health.status])
	}
	return result
}

func (cluster *Cluster) qourum() bool {
	result := 0
	for _, node := range cluster.nodes {
		if node.health.alive {
			result += 1
		}
	}
	return result > len(cluster.nodes)/2
}

type Health struct {
	status int
	alive  bool
}

type Node struct {
	addr      string
	msg_queue chan *WriteConsistencyMessage
	health *Health
}

func (node *Node) healthCheck() {
	for {
		_, err := http.Get(node.addr + "/ping")
		if err != nil {
			node.health.alive = false
			node.health.status = max(node.health.status-1, 0)
		} else {
			node.health.alive = true
			node.health.status = min(node.health.status+1, 2)
		}
		time.Sleep(3 * time.Second)
	}
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
