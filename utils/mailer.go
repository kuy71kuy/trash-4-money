package utils

import (
	"net/smtp"
	"os"
)

func NotifyPointEmail(point, pointBefore, pointAfter, receiverMail string) bool {

	senderMail := os.Getenv("MAILER_SENDER_MAIL")
	password := os.Getenv("MAILER_SENDER_PASS")

	mail := "From: " + senderMail + "\n" +
		"To: " + receiverMail + "\n" +
		"Subject: " + "Your Point Added" + "\n" +
		"Congratulation! You get " + point + " point!\n" +
		"Your point has increased from " + pointBefore + " to " + pointAfter

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", senderMail, password, "smtp.gmail.com"),
		senderMail, []string{receiverMail}, []byte(mail),
	)

	if err != nil {
		return false
	}
	return true
}
