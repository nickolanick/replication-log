package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"strconv"
)

func read_messages(w http.ResponseWriter, req *http.Request) {
	appDb := getAppDb()
	fmt.Fprintf(w, "read message\n")
}

func write_message(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "write message\n")
}

func main() {

	httpPort := getEnv("HTTP_PORT", "8080")
	role := getEnv("ROLE", "follower")
	delay, _ := strconv.Atoi(getEnv("DELAY", "10"))
	followers := strings.Split(os.Getenv("FOLLOWERS"), ",")

	fmt.Printf("%s", role)

	if len(followers) == 0 {
		panic("provide follower list")
	}

	fmt.Printf("%+v\n", getAppDb(role, followers, delay))

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
