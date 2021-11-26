package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// add worker pool
type AppDb struct {
	role         string
	followers    []string
	delay        int
	worker_queue []chan *WriteConsistencyMessage
	//should be hashmap and ordered array
	messages []string
}

// TODO: add lock/unlock
func (app *AppDb) write_message(m string) {
	if app.role == "leaders" {
		fmt.Printf("appending to secondaries")
	}
	app.messages = append(app.messages, m)
}

func (app AppDb) read_messages() []string {
	return app.messages
}

var appDb *AppDb
var lock = &sync.Mutex{}

func (appDb *AppDb) commitMessages(wcmsg *WriteConsistencyMessage) {
	// for chan in chans send message
	for _, queue := range appDb.worker_queue {
		queue <- wcmsg
	}
}

// TODO: consider waitgroup
func commitMessage(messages <-chan *WriteConsistencyMessage, follower string) {

	for wcmessage := range messages {
		// send request to commit
		postBody, _ := json.Marshal(map[string]string{
			"message": wcmessage.Message,
		})

		for {
			responseBody := bytes.NewBuffer(postBody)
			//Leverage Go's HTTP Post function to make request
			_, err := http.Post(follower+"/commit", "application/json", responseBody)
			// TODO: read response body status and retry
			if err == nil {
				break
			}
		}

    // negative wait group throws panic
    // we can catch it inside another function scope
    func (cond *sync.WaitGroup) {
        defer recover()
        cond.Done()
    }(wcmessage.WriteCond)
	}
}

func initAppDb(role string, followers []string, delay int) *AppDb {
	lock.Lock()
	defer lock.Unlock()

	if appDb == nil {
		fmt.Println("Creating single instance now.")
		appDb = &AppDb{role, followers, delay, []chan *WriteConsistencyMessage{}, []string{}}
	}
	// initialize buff chan
	// TODO: add to configuration length of queue
	for _, follower := range appDb.followers {
		queue := make(chan *WriteConsistencyMessage, 200000)
		go commitMessage(queue, follower)
		appDb.worker_queue = append(appDb.worker_queue, queue)
	}

	return appDb
}

func getAppDb() *AppDb {
	return appDb
}
