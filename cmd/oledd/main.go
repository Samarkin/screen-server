package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/samarkin/screen-server/auth"
	"github.com/samarkin/screen-server/engine"
)

const PASSWD_FILE_NAME = "./passwd"

// Health contains information about the server
type Health struct {
	OS           string `json:"os"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	h := Health{
		OS: runtime.GOOS,
	}
	e, err := engine.GetEngine()
	if e.Connected() {
		h.Status = "connected"
	} else {
		h.Status = "error"
		h.ErrorMessage = err.Error()
	}
	json.NewEncoder(w).Encode(h)
}

// MessageInfo contains information about a displayed message
type MessageInfo struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	e, _ := engine.GetEngine()
	var response [8]MessageInfo
	for i := 0; i < 8; i++ {
		response[i].Line = i
		response[i].Text = e.GetMessage(i)
	}
	json.NewEncoder(w).Encode(response)
}

// Message contains the text to display
type Message struct {
	Text     string `json:"text"`
	Duration *int   `json:"duration"`
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if msg.Duration != nil {
		http.Error(w, "Duration is not applicable here", http.StatusBadRequest)
		return
	}
	e, _ := engine.GetEngine()
	e.AppendMessage(msg.Text)
}

func handlePostPngImage(w http.ResponseWriter, r *http.Request) {
	e, _ := engine.GetEngine()
	var err error
	durationString := r.URL.Query().Get("duration")
	if len(durationString) > 0 {
		var duration int64
		if duration, err = strconv.ParseInt(durationString, 10, 32); err != nil {
			http.Error(w, "Invalid duration", http.StatusBadRequest)
			return
		}
		err = e.DisplayTemporaryImage(r.Body, time.Duration(duration)*time.Second)
	} else {
		err = e.DisplayImage(r.Body)
	}
	if err != nil {
		log.Printf("Unable to display image: %s", err)
		http.Error(w, "Unable to display the provided image", http.StatusBadRequest)
		return
	}
}

func handleDeleteMessages(w http.ResponseWriter, r *http.Request) {
	e, _ := engine.GetEngine()
	e.Clear()
}

func handlePutMessageOnLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	line64, err := strconv.ParseInt(vars["line"], 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	line := int(line64)
	decoder := json.NewDecoder(r.Body)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	e, _ := engine.GetEngine()
	if msg.Duration != nil {
		duration := *msg.Duration
		if duration > 3600 {
			duration = 3600
		} else if duration < 1 {
			duration = 1
		}
		e.DisplayTemporaryMessage(msg.Text, line, time.Duration(duration)*time.Second)
	} else {
		e.DisplayMessage(msg.Text, line)
	}
}

func handleGetMessageOnLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	line64, err := strconv.ParseInt(vars["line"], 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	line := int(line64)
	e, _ := engine.GetEngine()
	msg := MessageInfo{line, e.GetMessage(line)}
	json.NewEncoder(w).Encode(msg)
}

func handleDeleteMessageOnLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	line64, err := strconv.ParseInt(vars["line"], 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	line := int(line64)
	e, _ := engine.GetEngine()
	e.ClearMessage(line)
}

// LoginInfo contains login and password for authentication
type LoginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func handleLogin(context auth.AuthenticationContext, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginInfo LoginInfo
	if err := decoder.Decode(&loginInfo); err != nil {
		log.Printf("Authentication failed: %s", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
	} else if token, err := context.AuthenticateUser(loginInfo.Login, loginInfo.Password); err != nil {
		log.Printf("Authentication failed: %s", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
	} else {
		w.Header().Add("X-Session-Token", token)
	}
}

func loadPasswords(context auth.AuthenticationContext) {
	file, err := os.Open("./passwd")
	if err != nil {
		log.Println("Password file not found. Access to the main APIs will be disabled")
		log.Println("Use `go run github.com/samarkin/screen-server/cmd/create-user` to create users")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), ":")
		if len(columns) != 3 {
			continue
		}
		context.LoadUser(columns[0], columns[1], columns[2])
	}
}

func newRouter(loadPasswords func(auth.AuthenticationContext)) *mux.Router {
	r := mux.NewRouter()
	middleware, context := auth.NewAuthenticationMiddleware()
	r.Use(middleware)
	loadPasswords(context)
	context.ExcludeOperation("POST", "/api/login")
	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) { handleLogin(context, w, r) }).Methods("POST")
	r.HandleFunc("/api/health", handleGetHealth).Methods("GET")
	r.HandleFunc("/api/messages", handleGetMessages).Methods("GET")
	r.HandleFunc("/api/messages", handlePostMessage).Methods("POST")
	r.HandleFunc("/api/messages", handleDeleteMessages).Methods("DELETE")
	r.HandleFunc("/api/messages/{line:[0-7]}", handleGetMessageOnLine).Methods("GET")
	r.HandleFunc("/api/messages/{line:[0-7]}", handlePutMessageOnLine).Methods("PUT")
	r.HandleFunc("/api/messages/{line:[0-7]}", handleDeleteMessageOnLine).Methods("DELETE")
	r.HandleFunc("/api/image/png", handlePostPngImage).Methods("POST")
	return r
}

func main() {
	log.Printf("Initializing engine")
	e, _ := engine.GetEngine()
	defer e.Shutdown()
	r := newRouter(loadPasswords)
	server := &http.Server{
		Addr:    ":6533",
		Handler: r,
	}
	log.Printf("Listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
