package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

//chat-server
//user has to first login
//validation against user database
//once logged in list of contacts (online) is displayed
//user chooses to connect to user chosen
//display chat room
//connect them to eachother

type clientInfo struct {
	clientIP    net.IP
	connectWith string
}

type clientLogin struct {
	username string
	password string
}

type loginErrors struct {
	whatHappened string
}

func (l *clientLogin) register() {

}

func grantAccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chat.html")
}

func (l *loginErrors) Error() string {
	return l.whatHappened
}

func (l *clientLogin) lookupUser(w http.ResponseWriter) error {
	var userAndPassFields []string

	databaseHandler, _ := os.Open("users.txt")
	fileContents, err := os.ReadFile(databaseHandler.Name())

	if err != nil {
		log.Fatal(err)
	} else {
		var singleField []string
		//put usernames and passwords in a format thats easier to
		//work with

		for _, v := range fileContents {
			if v != '\n' {
				singleField = append(singleField, string(v))
			} else {
				joinedField := strings.Join(singleField, "")
				userAndPassFields = append(userAndPassFields, joinedField)
				singleField = singleField[:0]
			}
		}

		//find username and validate password
		for _, v := range userAndPassFields {
			for i, j := range v {
				if j == ':' {
					username := v[0:i]
					fmt.Println(username)

					if l.username == username {
						password := v[i:]

						if l.password == password {
							return nil
						} else {
							return &loginErrors{
								whatHappened: "Incorrect password.",
							}
						}
					}
				}
			}
		}

	}

	return &loginErrors{
		whatHappened: "Invalid username or password",
	}

}

func main() {

	loadLogin := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fileToOpen := r.URL.Path[1:]

			if string(fileToOpen) == "login.html" {
				http.ServeFile(w, r, string(fileToOpen))
			}

		case "POST":
			err := r.ParseForm()

			if err != nil {
				log.Fatal(err)
			} else {
				submittedLogin := &clientLogin{}

				username := r.FormValue("username")
				password := r.FormValue("password")

				submittedLogin.username = username
				submittedLogin.password = password

				err := submittedLogin.lookupUser(w)

				if err != nil {
					log.Fatal(err)
				} else {
					grantAccess(w, r)
				}

			}
		}
	}

	http.HandleFunc("/login.html", loadLogin)

	http.ListenAndServe(":8080", nil)
}
