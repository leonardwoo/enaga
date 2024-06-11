package enaga

import (
  "crypto/tls"
  "fmt"
  "github.com/emersion/go-imap"
  "github.com/emersion/go-imap/client"
  "github.com/emersion/go-message/mail"
  "github.com/emersion/go-sasl"
  "github.com/emersion/go-smtp"
  "log"
  "net/url"
  "strings"
)

/*
smtpServer is like smtp://host:port or smtps://host:port
*/
func Sender(smtpServer string, username string, password string, msg *MailMsg) (bool, error) {
  auth := sasl.NewPlainClient("", username, password)

  data := buildMailMsg(msg)

  urlObj, _ := url.Parse(smtpServer)
  smtpHostname := urlObj.Hostname() + ":" + urlObj.Port()
  if urlObj.Scheme == "smtp" {
    err := smtp.SendMail(smtpHostname, auth, msg.From, msg.To, strings.NewReader(data))
    if err != nil {
      return false, err
    }
  } else if urlObj.Scheme == "smtps" {
    err := smtp.SendMailTLS(smtpHostname, auth, msg.From, msg.To, strings.NewReader(data))
    if err != nil {
      return false, err
    }
  }

  return true, nil
}

/*
Receive with auto as read

imapServer is like imap://host:port or imaps://host:port
*/
func Receiver(imapServer string, username string, password string, enableSTARTTLS bool, mailSize uint32) (*MailMsg, error) {

  urlObj, _ := url.Parse(imapServer)
  imapHostname := urlObj.Hostname() + ":" + urlObj.Port()

  var c *client.Client
  var err error
  if urlObj.Scheme == "imap" {
    //c, err = client.Dial(imapHostname, nil)
    err = fmt.Errorf("Only supported IMAP with TLS")
  } else if urlObj.Scheme == "imaps" {
    c, err = client.DialTLS(imapHostname, nil)
  }

  if err != nil {
    return nil, err
  }

  if enableSTARTTLS {
    tlsConfig := &tls.Config{ServerName: urlObj.Hostname()}
    err = c.StartTLS(tlsConfig)
  }

  defer func(c *client.Client) {
    err = c.Logout()
  }(c)

  if err != nil {
    return nil, err
  }

  if err := c.Login(username, password); err != nil {
    return nil, fmt.Errorf("login failed %s", err.Error())
  }

  // Select INBOX
  mbox, err := c.Select("INBOX", false)
  if err != nil {
    return nil, err
  }

  // Get the last message
  if mbox.Messages == 0 {
    fmt.Println("No message in mailbox")
    return nil, nil
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
    err = c.Fetch(seqSet, items, messages)
  }()

  if err != nil {
    return nil, err
  }

  msg := <-messages
  if msg == nil {
    log.Println("server didn't returned message")
    return nil, nil
  }

  resp := msg.GetBody(&section)
  if resp == nil {
    log.Println("Server didn't returned message body")
    return nil, nil
  }

  // Create a new mail reader
  mr, err := mail.CreateReader(resp)
  if err != nil {
    return nil, err
  }

  return readMailMsg(mr), nil
}
