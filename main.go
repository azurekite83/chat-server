package main

import (
	"net/http"
)

//chat-server
//user has to first login
//validation against user database
//once logged in list of contacts (online) is displayed
//user chooses to connect to user chosen
//display chat room
//connect them to eachother

func main() {

	http.HandleFunc("/login.html", Login)
	http.HandleFunc("/register.html", Register)

	http.ListenAndServe(":8080", nil)
}
