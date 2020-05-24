package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samarkin/screen-server/engine"
)

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

func main() {
	log.Printf("Initializing engine")
	e := engine.GetEngine()
	defer e.Shutdown()

	r := mux.NewRouter()
	server := &http.Server{
		Addr:    ":6533",
		Handler: r,
	}
	r.HandleFunc("/api/health", handleGetHealth).Methods("GET")
	r.HandleFunc("/api/messages", handleGetMessages).Methods("GET")
	r.HandleFunc("/api/messages", handlePostMessage).Methods("POST")
	r.HandleFunc("/api/messages", handleDeleteMessages).Methods("DELETE")
	r.HandleFunc("/api/messages/{line:[0-7]}", handleGetMessageOnLine).Methods("GET")
	r.HandleFunc("/api/messages/{line:[0-7]}", handlePutMessageOnLine).Methods("PUT")
	r.HandleFunc("/api/messages/{line:[0-7]}", handleDeleteMessageOnLine).Methods("DELETE")

	log.Printf("Listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
