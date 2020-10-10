package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var r *mux.Router

func TestGetMessagesReturnsEmptyLines(t *testing.T) {
	r = newRouter()

	response := executeRequest("GET", "/api/messages", nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			assert.Equal(t, lines[i].Text, "", "Line %d", i)
		}
	})
}

func TestPutMessageChangesLine(t *testing.T) {
	r = newRouter()
	jsonStr := []byte(`{"text": "foobar"}`)
	response := executeRequest("PUT", "/api/messages/3", bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			if i != 3 {
				assert.Equal(t, lines[i].Text, "", "Line %d", i)
			} else {
				assert.Equal(t, lines[i].Text, "foobar", "Line %d", i)
			}
		}
	})
}

func TestDeleteMessagesClearsAll(t *testing.T) {
	r = newRouter()
	jsonStr := []byte(`{"text": "test"}`)
	response := executeRequest("PUT", "/api/messages/4", bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	response = executeRequest("DELETE", "/api/messages", nil)
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			assert.Equal(t, lines[i].Text, "", "Line %d", i)
		}
	})
}

func TestDeleteMessageClearsLine(t *testing.T) {
	r = newRouter()
	jsonStr := []byte(`{"text": "test4"}`)
	response := executeRequest("PUT", "/api/messages/4", bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	jsonStr = []byte(`{"text": "test5"}`)
	response = executeRequest("PUT", "/api/messages/5", bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	response = executeRequest("DELETE", "/api/messages/4", nil)
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			if i != 5 {
				assert.Equal(t, lines[i].Text, "", "Line %d", i)
			} else {
				assert.Equal(t, lines[i].Text, "test5", "Line %d", i)
			}
		}
	})
}

func executeRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertResponse(t *testing.T, response *httptest.ResponseRecorder, expectedCode int, expectedBody string) {
	assert.Equal(t, expectedCode, response.Code)
	body := response.Body.String()
	assert.Equal(t, strings.TrimSpace(body), expectedBody)
}

func assertMessageInfo(t *testing.T, response *httptest.ResponseRecorder, check func([]MessageInfo)) {
	assert.Equal(t, http.StatusOK, response.Code)
	decoder := json.NewDecoder(response.Body)
	var msg []MessageInfo
	if err := decoder.Decode(&msg); err != nil {
		t.Errorf("Response is not an array of MessageInfo's")
	}
	assert.Equal(t, 8, len(msg))
	check(msg)
}
