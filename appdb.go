package main

import (
	"fmt"
  "sync"
)

// add worker pool
type AppDb struct {
	role string
	followers []string
	delay int
	//should be hashmap and ordered array
	messages []string
}

// TODO: add lock/unlock
func (app AppDb) write_message(m string) {
	// for followers send shit
	if (app.role == "leaders") {
		fmt.Printf("appending to secondaries")
	}
	app.messages = append(app.messages, m)
}

func (app AppDb) read_messages() []string {
	return app.messages
}

var appDb *AppDb
var lock = &sync.Mutex{}

func getAppDb(role string, followers []string, delay int) *AppDb {
    if appDb == nil {
        lock.Lock()
        defer lock.Unlock()
        if appDb == nil {
            fmt.Println("Creating single instance now.")
            appDb = &AppDb{role, followers, delay, []string{}}
        } else {
            fmt.Println("Single instance already created.")
        }
    } else {
        fmt.Println("Single instance already created.")
    }

    return appDb
}
