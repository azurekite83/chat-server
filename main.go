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

func (l *loginErrors) Error() string {
	return l.whatHappened
}

func grantAccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chat.html")
}

func (l *clientLogin) getCredentials(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()

	if err != nil {
		return &loginErrors{"Could not parse request form."}
	} else {

		username := r.FormValue("username")
		password := r.FormValue("password")

		l.username = username
		l.password = password
	}

	return nil
}

func (l *clientLogin) lookupUser(onPage string) error {
	var userAndPassFields []string

	databaseHandler, _ := os.Open("users.txt")
	fileContents, err := os.ReadFile(databaseHandler.Name())

	if err != nil {
		log.Fatal(err)
	} else {
		closeErr, syncErr := databaseHandler.Close(), databaseHandler.Sync()

		if closeErr != nil {
			log.Fatal(closeErr)
		} else if syncErr != nil {
			log.Fatal(syncErr)
		}

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
		//if on login page
		if onPage == "login.html" {
			for _, v := range userAndPassFields {
				for i, j := range v {
					if j == ':' {
						username := v[0:i]

						if l.username == username {
							password := v[i+1:]

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

		if onPage == "register.html" {
			for _, v := range userAndPassFields {
				for i, j := range v {
					if j == ':' {
						username := v[0:i]

						if l.username == username {
							return &loginErrors{
								whatHappened: "Username taken.",
							}
						}
					}
				}
			}

			return nil
		}

	}

	return &loginErrors{
		whatHappened: "Invalid username or password",
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	pageToOpen := r.URL.Path[1:]

	switch r.Method {
	case "GET":
		if string(pageToOpen) == "login.html" {
			http.ServeFile(w, r, string(pageToOpen))
		}

	case "POST":
		submittedLogin := new(clientLogin)
		err := submittedLogin.getCredentials(w, r)

		if err != nil {
			log.Fatal(err)
		}

		err = submittedLogin.lookupUser(string(pageToOpen))

		if err != nil {
			log.Fatal(err)
		} else {
			grantAccess(w, r)
		}

	}
}

func (l *clientLogin) writeToDatabase() error {
	databaseHandle, fileError := os.OpenFile("users.txt", os.O_APPEND, 0600)

	if fileError != nil {
		log.Fatal(fileError)
	} else {
		usernameAndPass := l.username + ":" + l.password + "\n"
		fmt.Println(usernameAndPass)

		_, writeError := databaseHandle.WriteString(usernameAndPass)

		if writeError != nil {
			return writeError
		}
	}

	return nil
}

func register(w http.ResponseWriter, r *http.Request) {
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
		submittedRegistration := new(clientLogin)
		parseError := submittedRegistration.getCredentials(w, r)

		if parseError != nil {
			log.Fatal(parseError)
		}

		userExists := submittedRegistration.lookupUser(onPage)

		if userExists != nil {
			log.Fatal(userExists)
		} else {
			writeError := submittedRegistration.writeToDatabase()

			if writeError != nil {
				log.Fatal(writeError)
			}
		}
	}

}

func main() {

	http.HandleFunc("/login.html", login)
	http.HandleFunc("/register.html", register)

	http.ListenAndServe(":8080", nil)
}
