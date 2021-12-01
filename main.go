package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func read_messages(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", repository.GetMessages())
}

func write_message(w http.ResponseWriter, req *http.Request) {
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

	if config.role == "follower" {
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
    wr_cons_msg.TotalOrder = counter.get()
    fmt.Println(wr_cons_msg.TotalOrder)

		cluster.commitMessages(&wr_cons_msg)
		wr_cons_msg.WriteCond.Wait()
	}

	fmt.Fprintf(w, "write consistency %s\n", wr_cons_msg)
}

func commit_message(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Duration(config.delay) * time.Second)

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

	repository.AppendMessage(wr_cons_msg)
	return
}

func main() {

	// move this to read_config
	config.init()
	cluster.init(config.nodes)

	fmt.Printf("%s", config.role)

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

	http.ListenAndServe(":"+config.httpPort, nil)
}
