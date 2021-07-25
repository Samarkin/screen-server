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

	"github.com/gorilla/mux"
	"github.com/samarkin/screen-server/auth"
	"github.com/samarkin/screen-server/engine"
)

const PASSWD_FILE_NAME = "./passwd"

// Health contains information about the server
type Health struct {
	OS     string `json:"os"`
	Status string `json:"status"`
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	h := Health{
		OS: runtime.GOOS,
	}
	if engine.GetEngine().Connected() {
		h.Status = "connected"
	} else {
		h.Status = "error"
	}
	json.NewEncoder(w).Encode(h)
}

// MessageInfo contains information about a displayed message
type MessageInfo struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	var response [8]MessageInfo
	for i := 0; i < 8; i++ {
		response[i].Line = i
		response[i].Text = engine.GetEngine().GetMessage(i)
	}
	json.NewEncoder(w).Encode(response)
}

// Message contains the text to display
type Message struct {
	Text string `json:"text"`
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	engine.GetEngine().AppendMessage(msg.Text)
}

func handleDeleteMessages(w http.ResponseWriter, r *http.Request) {
	engine.GetEngine().Clear()
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
	engine.GetEngine().DisplayMessage(msg.Text, line)
}

func handleGetMessageOnLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	line64, err := strconv.ParseInt(vars["line"], 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	line := int(line64)
	msg := MessageInfo{line, engine.GetEngine().GetMessage(line)}
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
	engine.GetEngine().ClearMessage(line)
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
	return r
}

func main() {
	log.Printf("Initializing engine")
	e := engine.GetEngine()
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
