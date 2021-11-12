package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Message          string `json:"message"`
	WriteConsistency int64  `json:"write_consistency"`
}

func read_messages(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	fmt.Fprintf(w, "read message %s\n", appDb.messages)
}

func write_message(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	var m Message

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if m.WriteConsistency == 0 {
		m.WriteConsistency = int64(len(appDb.followers))
	}
	if appDb.role == "follower" {
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(m)
		resp, err := http.Post("http://leader:5000", "application/json", payloadBuf)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	} else {
		wcmsg := WriteConsistencyMessage{m.Message, m.WriteConsistency}
		appDb.commitMessages(&wcmsg)

		for {
			if wcmsg.wc_counter <= 0 {
				break
			}
		}
	}
	fmt.Fprintf(w, "write message %s write consist %s\n", appDb.followers, m.WriteConsistency)
}

func commit_message(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	time.Sleep(time.Duration(appDb.delay) * time.Second)

	var m Message
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "message: %+v", m)
	appDb.write_message(m.Message)
	return
}

func main() {

	httpPort := getEnv("HTTP_PORT", "8080")
	role := getEnv("ROLE", "follower")
	delay, _ := strconv.Atoi(getEnv("DELAY", "0"))
	followers := strings.Split(os.Getenv("FOLLOWERS"), ",")

	fmt.Printf("%s", role)

	if len(followers) == 0 {
		panic("provide follower list")
	}

	fmt.Printf("%+v\n", initAppDb(role, followers, delay))

	http.HandleFunc("/commit", commit_message)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			read_messages(w, req)
		case http.MethodPost:
			write_message(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+httpPort, nil)
}
