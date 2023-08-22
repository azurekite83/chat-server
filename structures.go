package main

import (
	"net"
)

type ClientInfo struct {
	clientIP net.IP
	userID   string
	password string
}

type LoginErrors struct {
	whatHappened string
}
