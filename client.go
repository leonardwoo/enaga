package enaga

import (
  "crypto/tls"
  "fmt"
  "github.com/emersion/go-imap"
  "github.com/emersion/go-imap/client"
  "github.com/emersion/go-message/mail"
  "github.com/emersion/go-sasl"
  "github.com/emersion/go-smtp"
  "net/url"
  "strings"
)

/*
smtpServer is like smtp://host:port or smtps://host:port
*/
func Sender(smtpServer string, username string, password string, msg *MailMsg) {
  auth := sasl.NewPlainClient("", username, password)

  data := buildMailMsg(msg)

  urlObj, _ := url.Parse(smtpServer)
  smtpHostname := urlObj.Hostname() + ":" + urlObj.Port()
  if urlObj.Scheme == "smtp" {
    err := smtp.SendMail(smtpHostname, auth, msg.From, msg.To, strings.NewReader(data))
    if err != nil {
      panic(err)
    }
  } else if urlObj.Scheme == "smtps" {
    err := smtp.SendMailTLS(smtpHostname, auth, msg.From, msg.To, strings.NewReader(data))
    if err != nil {
      panic(err)
    }
  }
}

/*
Receive with auto as read

imapServer is like imap://host:port or imaps://host:port
*/
func Receiver(imapServer string, username string, password string, enableSTARTTLS bool, mailSize uint32) {

  urlObj, _ := url.Parse(imapServer)
  imapHostname := urlObj.Hostname() + ":" + urlObj.Port()

  var c *client.Client
  var err error
  if urlObj.Scheme == "imap" {
    //c, err = client.Dial(imapHostname, nil)
    panic("Only supported IMAP with TLS")
  } else if urlObj.Scheme == "imaps" {
    c, err = client.DialTLS(imapHostname, nil)
  }

  if err != nil {
    panic(err)
  }

  if enableSTARTTLS {
    tlsConfig := &tls.Config{ServerName: urlObj.Hostname()}
    err = c.StartTLS(tlsConfig)
  }

  defer func(c *client.Client) {
    err := c.Logout()
    if err != nil {
      panic(err)
    }
  }(c)

  if err := c.Login(username, password); err != nil {
    panic("login failed " + err.Error())
  }

  // Select INBOX
  mbox, err := c.Select("INBOX", false)
  if err != nil {
    panic(err)
  }

  // Get the last message
  if mbox.Messages == 0 {
    fmt.Println("No message in mailbox")
    return
  }

  from := mailSize
  to := mbox.Messages
  if mbox.Messages > mailSize {
    // We're using unsigned integers here, only substract if the result is > 0
    from = mbox.Messages - mailSize
  }
  seqSet := new(imap.SeqSet)
  seqSet.AddRange(from, to)

  // Get the whole message body
  var section imap.BodySectionName
  items := []imap.FetchItem{section.FetchItem()}

  messages := make(chan *imap.Message, mailSize)
  go func() {
    if err := c.Fetch(seqSet, items, messages); err != nil {
      panic(err)
    }
  }()

  msg := <-messages
  if msg == nil {
    fmt.Println("server didn't returned message")
    return
  }

  resp := msg.GetBody(&section)
  if resp == nil {
    fmt.Println("Server didn't returned message body")
    return
  }

  // Create a new mail reader
  mr, err := mail.CreateReader(resp)
  if err != nil {
    panic(err)
  }

  mailMsg := readMailMsg(mr)
  fmt.Println(toJsonString(mailMsg))
}
