package enaga

import "testing"

func TestSender(t *testing.T) {
  mailMsg := &MailMsg{
    From:    "test@matoro.net",
    To:      []string{"leo@matoro.net"},
    Subject: "Go mail test",
    Body:    "gomail test with emersion go-smtp",
  }

  //Sender("smtp://smtp.ym.163.com:25", "test@matoro.net", "3wJzYcy827", mailMsg)
  Sender("smtps://smtp.ym.163.com:994", "test@matoro.net", "3wJzYcy827", mailMsg)
}

func TestReceiver(t *testing.T) {
  //Receiver("imap://imap.ym.163.com:143", "test@matoro.net", "3wJzYcy827", false, 20)
  Receiver("imaps://imap.ym.163.com:993", "test@matoro.net", "3wJzYcy827", true, 20)
}
