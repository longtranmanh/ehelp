package utils

import "gopkg.in/gomail.v2"

type Mail struct {
	Subject string
	Body    string
	To      string
}

var mailDialer = gomail.NewDialer("smtp.gmail.com", 465, "trunglenlvn@gmail.com", "gfgsimshbzgwrxwa")

func (mail Mail) Send() {
	m := gomail.NewMessage()
	m.SetHeader("From", "trunglenlvn@gmail.com")
	m.SetHeader("To", mail.To)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)
	// Send the email to Bob, Cora and Dan.
	if err := mailDialer.DialAndSend(m); err != nil {
		panic(err)
	}
}
