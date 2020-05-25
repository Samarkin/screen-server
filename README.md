# ScreenServer

The goal of this project is to develop a remotely controlled OLED screen.

## Technology Stack
* Raspberry Pi 3 Model B (to be potentially replaced with a cheaper alternative)
* [SH1106](https://www.displayfuture.com/Display/datasheet/controller/SH1106.pdf)-driven 128x64 OLED screen
* [Golang](https://golang.org)

## To Build
1. Get the source code
```shell
go get github.com/samarkin/screen-server/oledd
```
2. Build and run
```shell
go run github.com/samarkin/screen-server/oledd
```

## Sample Usage
1. Find out IP of your Raspberry Pi. For example, `192.168.1.5`.
2. From any computer connected to the same network make an HTTP request to port 6533
(OLED spelled on [phone keypad](https://en.wikipedia.org/wiki/E.161))
```shell
curl -v --data '{"text": "Hello, world!"}' http://192.168.1.5:6533/api/messages
```

[Full API Description](oledd/API.md)
