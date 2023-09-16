package main

import (
	"net"
)

type ClientInfo struct {
	IPAddr     string
	connStatus *net.Conn
	userID     string
	password   string
}

type LoginErrors struct {
	whatHappened string
}

var ClientPool = make(map[string]ClientInfo)
