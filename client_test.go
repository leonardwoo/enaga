package enaga

import (
  "fmt"
  "testing"
)

func TestSender(t *testing.T) {
  mailMsg := &MailMsg{
    From:    "",
    To:      []string{""},
    Subject: "",
    Body:    "",
  }

  //Sender("smtp://smtp.example.com:25", "", "", mailMsg)
  flag, err := Sender("smtps://smtp.example.com:994", "", "", mailMsg)
  if err != nil {
    panic(err)
  }
  if flag {
    fmt.Println("sender successful")
  } else {
    fmt.Println("sender failed")
  }
}

func TestReceiver(t *testing.T) {
  //Receiver("imap://imap.example.com:143", "", "", false, 20)
  mails, err := Receiver("imaps://imap.example.com:993", "", "", true, 20)
  if err != nil {
    panic(err)
  }

  if mails != nil {
    fmt.Println(toJsonString(mails))
  }
}
