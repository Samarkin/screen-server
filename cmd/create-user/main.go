package main

import (
	"fmt"
	"os"

	"github.com/samarkin/screen-server/auth"
	"golang.org/x/term"
)

const PASSWD_FILE_NAME = "./passwd"

func main() {
	fmt.Print("Enter new login: ")
	var login string
	if _, err := fmt.Scanln(&login); err != nil {
		fmt.Printf("An error occured while reading input: %s\n", err)
		os.Exit(1)
	}

	fmt.Print("Enter password: ")
	var password string
	if pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd())); err != nil {
		fmt.Printf("An error occured while reading input: %s\n", err)
		os.Exit(1)
	} else {
		password = string(pwdBytes)
	}
	fmt.Println()
	pwdInfo := auth.HashPassword(password)
	if file, err := os.OpenFile(PASSWD_FILE_NAME, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		fmt.Printf("Failed to open %s: %s", PASSWD_FILE_NAME, err)
		os.Exit(2)
	} else {
		defer file.Close()
		fmt.Fprintf(file, "%s:%s:%s\n", login, pwdInfo.Salt, pwdInfo.Hash)
	}
}
