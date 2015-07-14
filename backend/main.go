package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

var emailList = []string{"antoine.pourchet@gmail.com"}
var lastCount = -1

func removeEmail(toremove string) []string {
    newEmailList := []string{}
    for _, email := range emailList {
		if email != toremove {
			newEmailList = append(newEmailList, email)
		}
	}
	return newEmailList
}

func cleanBody(body []byte) string {
	return strings.Replace(string(body), "\n", "", -1)
}

func sendEmail(count int) {
	fmt.Println("Sending the email")

	subjectString := "ALERT: INVENTORY LOW"
	bodyString := ""
	if count == 0 {
		bodyString = "Out of items, please refill"
	} else if count == 1 {
		bodyString = "You only have 1 item left, order some moar"
	} else if count < 3 {
		bodyString = "Running low, be aware"
	}

	emailUser := &EmailUser{"monitor.inventory", "squaresquaresquare1!", "smtp.gmail.com", 587}
	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer)

	emailBody := fmt.Sprintf("From: AlwaysBeer\nTo: Dear Customer\nSubject: %s\n\n%s\nBottles Left: %d\n", subjectString, bodyString, count)
	err := smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), auth,
		emailUser.Username,
		emailList,
		[]byte(emailBody))
	if err != nil {
		fmt.Println("ERROR: attempting to send a mail ", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: statusHandler")
	fmt.Fprintf(w, "ok")
	w.WriteHeader(200)
}

func addEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: addEmailHandler")
	email := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		email = cleanBody(body)
	} else {
		email = r.URL.Query().Get("email")
	}
	fmt.Printf("Email: '%s'\n", email)
	if len(email) != 0 {
		emailList = append(emailList, email)
	}
	fmt.Printf("Email List: %v\n\n", emailList)
	fmt.Fprintf(w, "Email added: %s", email)
}

func lastCountHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: lastCountHandler")
	if lastCount < 0 {
		fmt.Fprintf(w, "We do not know how many bottles are left!")
	} else {
		fmt.Fprintf(w, "Last Count: %d\n", lastCount)
	}
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: countHandler")
	count := 0
	if len(r.URL.RawQuery) == 0 {
		bodyArr, _ := ioutil.ReadAll(r.Body)
		body := cleanBody(bodyArr)
		count, _ = strconv.Atoi(body)
	} else {
		countStr := r.URL.Query().Get("count")
		count, _ = strconv.Atoi(countStr)
	}

	fmt.Printf("Got a bottle count: %d\n", count)
	fmt.Fprintf(w, "Bottle count: '%d'\n", count)
	if count < 3 {
		go sendEmail(count)
	}
	lastCount = count
}

func removeEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: removeEmailHandler")
	email := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		email = cleanBody(body)
	} else {
		email = r.URL.Query().Get("email")
	}
	fmt.Printf("Email: '%s'\n", email)
	if len(email) > 0 {
		emailList = removeEmail(email)
	}
	fmt.Printf("Email List: %v\n\n", emailList)
	fmt.Fprintf(w, "Email removed: %s", email)
}

func main() {
	http.HandleFunc("/_status", statusHandler)
	http.HandleFunc("/bottlecount", countHandler)
	http.HandleFunc("/addemail", addEmailHandler)
	http.HandleFunc("/lastcount", lastCountHandler)
	http.HandleFunc("/removeemail", removeEmailHandler)
	http.ListenAndServe(":8080", nil)
}
