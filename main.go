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
  "sync"
)

func read_messages(w http.ResponseWriter, req *http.Request) {
	//appDb := getAppDb()
  fmt.Fprintf(w, "%s\n", repository.GetMessages())
  //fmt.Fprintf(w, "%s\n", appDb.messages)
}

func write_message(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	var wr_cons_msg WriteConsistencyMessage

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&wr_cons_msg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if appDb.role == "follower" {
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(wr_cons_msg)
		resp, err := http.Post("http://leader:5000", "application/json", payloadBuf)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	} else {
		// if role is leader we send to channel commit all
		// otherwise we proxy to leader
		// send commitMessage
		// wc should be from field, default field == follower number
		// while wcmsg atomic counter >= 0 wait

		if wr_cons_msg.WriteConsistency == 0 {
			wr_cons_msg.WriteConsistency = len(appDb.followers)
		}

    var write_cond sync.WaitGroup
    write_cond.Add(wr_cons_msg.WriteConsistency)

    wr_cons_msg.WriteCond = &write_cond

		appDb.commitMessages(&wr_cons_msg)

    write_cond.Wait()
	}

	fmt.Fprintf(w, "write message %s write consist %i\n", appDb.followers, wr_cons_msg.WriteConsistency)
}

func commit_message(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	time.Sleep(time.Duration(appDb.delay) * time.Second)

	var wr_cons_msg WriteConsistencyMessage
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&wr_cons_msg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "message: %+v", wr_cons_msg)
	repository.AppendMessage(wr_cons_msg.Message)
	return
}

func main() {

  // move this to read_config
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
