package main

import (
	"log"
	"net/http"
)

func grantAccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chat.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	pageToOpen := r.URL.Path[1:]

	switch r.Method {
	case "GET":
		if string(pageToOpen) == "login.html" {
			http.ServeFile(w, r, string(pageToOpen))
		}

	case "POST":
		submittedLogin := new(ClientLogin)
		err := submittedLogin.GetCredentials(w, r)

		if err != nil {
			log.Fatal(err)
		}

		err = submittedLogin.LookupUser(string(pageToOpen))

		if err != nil {
			log.Fatal(err)
		} else {
			grantAccess(w, r)
		}

	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	//This is a bit of obscure code, the error here is from the
	//lookupUser function, which if you look at there's
	//functionality to see what page the user is on.
	//
	//Since the code from the registration page executes,
	//either an error explaining that the username is taken
	//is returned or nothing at all

	onPage := r.URL.Path[1:]

	if r.Method == "GET" {
		http.ServeFile(w, r, "register.html")
	}
	if r.Method == "POST" {
		submittedRegistration := new(ClientLogin)
		parseError := submittedRegistration.GetCredentials(w, r)

		if parseError != nil {
			log.Fatal(parseError)
		}

		userExists := submittedRegistration.LookupUser(onPage)

		if userExists != nil {
			log.Fatal(userExists)
		} else {
			writeError := submittedRegistration.WriteToDatabase()

			if writeError != nil {
				log.Fatal(writeError)
			}
		}
	}

}
