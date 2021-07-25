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
	"github.com/samarkin/screen-server/auth"
	"github.com/stretchr/testify/assert"
)

var r *mux.Router

func TestAuthenticationRequired(t *testing.T) {
	r = newRouter(createFakeUser)
	response := executeRequest("GET", "/api/messages", "", nil)
	assertResponse(t, response, http.StatusUnauthorized, "Unauthorized")
}

func TestGetMessagesReturnsEmptyLines(t *testing.T) {
	r = newRouter(createFakeUser)
	token := login(t)
	response := executeRequest("GET", "/api/messages", token, nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			assert.Equal(t, lines[i].Text, "", "Line %d", i)
		}
	})
}

func TestPutMessageChangesLine(t *testing.T) {
	r = newRouter(createFakeUser)
	token := login(t)
	jsonStr := []byte(`{"text": "foobar"}`)
	response := executeRequest("PUT", "/api/messages/3", token, bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", token, nil)

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
	r = newRouter(createFakeUser)
	token := login(t)
	jsonStr := []byte(`{"text": "test"}`)
	response := executeRequest("PUT", "/api/messages/4", token, bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	response = executeRequest("DELETE", "/api/messages", token, nil)
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", token, nil)

	assertMessageInfo(t, response, func(lines []MessageInfo) {
		for i := range lines {
			assert.Equal(t, lines[i].Line, i)
			assert.Equal(t, lines[i].Text, "", "Line %d", i)
		}
	})
}

func TestDeleteMessageClearsLine(t *testing.T) {
	r = newRouter(createFakeUser)
	token := login(t)
	jsonStr := []byte(`{"text": "test4"}`)
	response := executeRequest("PUT", "/api/messages/4", token, bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	jsonStr = []byte(`{"text": "test5"}`)
	response = executeRequest("PUT", "/api/messages/5", token, bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusOK, "")
	response = executeRequest("DELETE", "/api/messages/4", token, nil)
	assertResponse(t, response, http.StatusOK, "")

	response = executeRequest("GET", "/api/messages", token, nil)

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

func TestJsonRequiredWhenLoggingIn(t *testing.T) {
	r = newRouter(createFakeUser)
	nonJsonStr := []byte(`login=admin&password=admin`)
	response := executeRequest("POST", "/api/login", "", bytes.NewBuffer(nonJsonStr))
	assertResponse(t, response, http.StatusForbidden, "Forbidden")
}

func TestPasswordIsVerifiedWhenLoggingIn(t *testing.T) {
	r = newRouter(createFakeUser)
	jsonStr := []byte(`{"login": "admin", "password": "invalid_password"}`)
	response := executeRequest("POST", "/api/login", "", bytes.NewBuffer(jsonStr))
	assertResponse(t, response, http.StatusForbidden, "Forbidden")
}

func createFakeUser(context auth.AuthenticationContext) {
	hash := auth.HashPassword("admin")
	context.LoadUser("admin", hash.Salt, hash.Hash)
}

func login(t *testing.T) string {
	jsonStr := []byte(`{"login": "admin", "password": "admin"}`)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	token := rr.Header().Get("X-Session-Token")
	return token
}

func executeRequest(method, url, token string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, body)
	if token != "" {
		req.Header.Add("X-Session-Token", token)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertResponse(t *testing.T, response *httptest.ResponseRecorder, expectedCode int, expectedBody string) {
	body := response.Body.String()
	if assert.Equal(t, expectedCode, response.Code, "Body: %s", body) {
		assert.Equal(t, strings.TrimSpace(body), expectedBody)
	}
}

func assertMessageInfo(t *testing.T, response *httptest.ResponseRecorder, check func([]MessageInfo)) {
	if assert.Equal(t, http.StatusOK, response.Code) {
		decoder := json.NewDecoder(response.Body)
		var msg []MessageInfo
		if err := decoder.Decode(&msg); err != nil {
			t.Errorf("Response is not an array of MessageInfo's")
		}
		assert.Equal(t, 8, len(msg))
		check(msg)
	}
}
