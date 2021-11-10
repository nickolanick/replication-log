package main

import (
	"fmt"
	"net/http"
	"sync"
)

// add worker pool
type AppDb struct {
	role      string
	followers []string
	delay     int
	//should be hashmap and ordered array
	messages []string
	// we use this to handle termination -> no idea what we actually handle but let it be
	wg sync.WaitGroup
}

// TODO: add lock/unlock
func (app *AppDb) write_message(m string) {
	// for followers send shit
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

func commitMessage(messages chan string, wg *sync.WaitGroup, follower string) {
	defer wg.Done()

	for message := range messages {
		// send request to follower here
	}
}

func initAppDb(role string, followers []string, delay int) *AppDb {
	lock.Lock()
	defer lock.Unlock()
	if appDb == nil {
		fmt.Println("Creating single instance now.")
		appDb = &AppDb{role, followers, delay, []string{}}
	}
	return appDb
}

func getAppDb() *AppDb {
	return appDb
}
