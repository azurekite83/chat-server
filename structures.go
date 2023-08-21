package main

import (
	"net"
)

type ClientInfo struct {
	clientIP net.IP
	userID   string
}

type ClientLogin struct {
	username string
	password string
}

type LoginErrors struct {
	whatHappened string
}
