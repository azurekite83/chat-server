package main

import (
	"net/http"
)

//chat-server
//user has to first login
//validation against user database
//once logged in list of contacts (online) is displayed
//
//
//Change server functionality to listen on a port,
//return connection struct to whatever function handling clients,
//respond with index page. On login, gather all clients that havent
// timed out, display them as online and then go from there
//
//
//
//user chooses to connect to user chosen
//display chat room
//connect them to eachother

func main() {

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/login.html", Login)
	http.HandleFunc("/register.html", Register)

	server := http.Server{
		Addr:      "127.0.0.1:8080",
		ConnState: ConnStateHandler,
	}

	server.ListenAndServe()
}

//TODO: So much to do. Kill me.
//ConnStateHandler is not behaving like its supposed to
//request handlers need to be directed to the right file.
//Need to better understand how tf connections work
//implement tests, this fmt.Println shit is not gonna cut it
