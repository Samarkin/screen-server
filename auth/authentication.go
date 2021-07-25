package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/sha3"
)

// AuthenticationContext holds the context of all authentication operations
type AuthenticationContext interface {
	AuthenticateUser(login, password string) (string, error)
	ExcludeOperation(method, uri string)
	LoadUser(login, salt, hash string) error
}

// PasswordInformation consists of salt and hash
type PasswordInformation struct {
	Salt string
	Hash string
}

type authenticationContext struct {
	excludedOperations     map[string]bool
	authenticationSessions map[string]string
	registeredUsers        map[string]PasswordInformation
}

type authenticationMiddleware struct {
	next    http.Handler
	context *authenticationContext
}

// NewAuthenticationMiddleware creates a new authentication context and provides a middleware func to build middleware instances
func NewAuthenticationMiddleware() (mux.MiddlewareFunc, AuthenticationContext) {
	context := &authenticationContext{make(map[string]bool), make(map[string]string), make(map[string]PasswordInformation)}
	middlewareFunc := func(next http.Handler) http.Handler {
		return authenticationMiddleware{next, context}
	}
	return middlewareFunc, context
}

func (amw authenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Session-Token")
	log.Printf(r.Method + " " + r.RequestURI)
	if user, found := amw.context.authenticationSessions[token]; found {
		log.Printf("Authenticated user %s\n", user)
		amw.next.ServeHTTP(w, r)
	} else if _, found := amw.context.excludedOperations[r.Method+" "+r.RequestURI]; found {
		log.Printf("No auth required")
		amw.next.ServeHTTP(w, r)
	} else {
		log.Printf("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func hashPassword(password, salt string) string {
	hashBytes := make([]byte, 64)
	sha3.ShakeSum256(hashBytes, []byte(password+salt))
	return hex.EncodeToString(hashBytes)
}

// HashPassword generates random salt and calculates password hash
func HashPassword(password string) PasswordInformation {
	saltBytes := make([]byte, 4)
	rand.Read(saltBytes)
	salt := hex.EncodeToString(saltBytes)
	hash := hashPassword(password, salt)
	return PasswordInformation{salt, hash}
}

// AuthenticateUser authenticates user and returns a session token
func (context *authenticationContext) AuthenticateUser(login, password string) (string, error) {
	if passwordInfo, found := context.registeredUsers[login]; found {
		if hashPassword(password, passwordInfo.Salt) == passwordInfo.Hash {
			token := uuid.NewString()
			context.authenticationSessions[token] = login
			return token, nil
		}
	}
	return "", errors.New("invalid username/password combination")
}

// ExcludeOperation adds an operation to the list of operations that don't require authentication
func (context *authenticationContext) ExcludeOperation(method, uri string) {
	context.excludedOperations[method+" "+uri] = true
}

// LoadUser adds a new user given their login, salt and hash information
func (context *authenticationContext) LoadUser(login, salt, hash string) error {
	if _, found := context.registeredUsers[login]; found {
		return errors.New("user already exists")
	}
	context.registeredUsers[login] = PasswordInformation{salt, hash}
	return nil
}
