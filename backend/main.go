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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
	w.WriteHeader(200)
}

func addEmailHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	emailList = append(emailList, cleanBody(body))
	fmt.Println(emailList)
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

	emailUser := &EmailUser{"alwaysbeer.gypsies", "squaresquaresquare!", "smtp.gmail.com", 587}
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

func countHandler(w http.ResponseWriter, r *http.Request) {
	bodyArr, _ := ioutil.ReadAll(r.Body)
	body := cleanBody(bodyArr)
	count, _ := strconv.Atoi(body)

	fmt.Fprintf(w, "Bottle count: '%d'\n", count)
	if count < 3 {
		go sendEmail(count)
	}
}

func main() {
	http.HandleFunc("/_status", statusHandler)
	http.HandleFunc("/bottlecount", countHandler)
	http.HandleFunc("/addemail", addEmailHandler)
	http.ListenAndServe(":8080", nil)
}
