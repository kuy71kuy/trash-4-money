package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func NotifyPointEmail(point, pointBefore, pointAfter, receiverMail, name string) bool {

	senderMail := os.Getenv("MAILER_SENDER_MAIL")
	password := os.Getenv("MAILER_SENDER_PASS")
	fmt.Println()
	mail := "From: " + senderMail + "\n" +
		"To: " + receiverMail + "\n" +
		"Subject: " + "Your Point Added" + "\n" +
		"Congratulation " + name + "! You get " + point + " point!\n" +
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

func NotifyPaymentEmail(point, numberType, receiverMail, name string) bool {

	senderMail := os.Getenv("MAILER_SENDER_MAIL")
	password := os.Getenv("MAILER_SENDER_PASS")
	fmt.Println()
	mail := "From: " + senderMail + "\n" +
		"To: " + receiverMail + "\n" +
		"Subject: " + "Your Payment Done" + "\n" +
		"Congratulation " + name + "! Your payment of " + point + " has done\n" +
		"Please check your balance in your " + numberType

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", senderMail, password, "smtp.gmail.com"),
		senderMail, []string{receiverMail}, []byte(mail),
	)

	if err != nil {
		return false
	}
	return true
}
