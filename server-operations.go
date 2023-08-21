package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func (l *ClientLogin) WriteToDatabase() error {
	databaseHandle, fileError := os.OpenFile("users.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer databaseHandle.Close()

	if fileError != nil {
		log.Fatal(fileError)
	} else {
		usernameAndPass := l.username + ":" + l.password + "\n"

		_, writeError := databaseHandle.WriteString(usernameAndPass)

		if writeError != nil {
			return writeError
		}

		if closeErr := databaseHandle.Close(); closeErr != nil {
			return closeErr
		}
	}

	return nil
}

func (l *ClientLogin) LookupUser(onPage string) error {
	var userAndPassFields []string

	fileContents, err := os.ReadFile("users.txt")

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
								return &LoginErrors{
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
							return &LoginErrors{
								whatHappened: "Username taken.",
							}
						}
					}
				}
			}

			return nil
		}

	}

	return &LoginErrors{
		whatHappened: "Invalid username or password",
	}

}

func (l *LoginErrors) Error() string {
	return l.whatHappened
}

func (l *ClientLogin) GetCredentials(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()

	if err != nil {
		return &LoginErrors{"Could not parse request form."}
	} else {

		username := r.FormValue("username")
		password := r.FormValue("password")

		l.username = username
		l.password = password
	}

	return nil
}
