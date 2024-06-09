package enaga

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "github.com/emersion/go-message/mail"
  "io"
  "regexp"
  "strings"
  "time"
)

type MailMsg struct {
  From      string
  To        []string
  Cc        []string
  Bcc       []string
  ReplyTo   []string
  Subject   string
  Body      string
  IsHtml    bool `default:"false"`
  Timestamp time.Time
  MessageId string
}

//func (addr MailAddr) mailAddrToString() string {
//  if len(addr.displayName) == 0 {
//    return addr.mailAddr
//  }
//  return addr.displayName + " <" + addr.mailAddr + ">"
//}
//
//func (addr MailAddr) mailAddrToStringWithUtf8() string {
// if len(addr.displayName) == 0 {
//   return addr.mailAddr
// }
// return toUtf8(addr.displayName) + " <" + addr.mailAddr + ">"
//}

func mailAddrToString(addr *mail.Address) string {
  if len(addr.Name) == 0 {
    return addr.Address
  }
  return addr.Name + " <" + addr.Address + ">"
}

func mailAddrsToStrings(addrs []*mail.Address) []string {
  arrStrs := make([]string, len(addrs))
  for i := 0; i < len(addrs); i++ {
    arrStrs[i] = mailAddrToString(addrs[i])
  }
  return arrStrs
}

func mailAddrsToString(addrs []*mail.Address) string {
  var addrsStr = ""
  for i := 0; i < len(addrs); i++ {
    addrsStr += mailAddrToString(addrs[i]) + ", "
  }
  return addrsStr[0 : len(addrsStr)-2]
}

func toUtf8(text string) string {
  base64Str := base64.StdEncoding.EncodeToString([]byte(text))
  return "=?UTF-8?B?" + base64Str + "?="
}

func fromUtf8(utf8Txt string) string {
  base64Str := utf8Txt[10 : len(utf8Txt)-2]
  raw, err := base64.StdEncoding.DecodeString(base64Str)
  if err != nil {
    panic(err)
  }
  return string(raw)
}

func checkAscii(text string) bool {
  matched, err := regexp.MatchString("^[\x20-\x7E]+$", text)
  if err != nil {
    return false
  }
  return matched
}

func toJsonString(v any) (string, error) {
  b, err := json.Marshal(v)
  return string(b), err
}

func buildMailMsg(msg *MailMsg) string {
  var data string

  to := strings.Join(msg.To, ",\r\n ")
  data += "To: " + to + "\r\n"

  if msg.Cc != nil {
    cc := strings.Join(msg.Cc, ",\r\n ")
    data += "Cc: " + cc + "\r\n"
  }

  if msg.Bcc != nil {
    bcc := strings.Join(msg.Bcc, ",\r\n ")
    data += "Bcc: " + bcc + "\r\n"
  }

  if msg.ReplyTo != nil {
    replyTo := strings.Join(msg.ReplyTo, ",\r\n ")
    data += "Reply-To: " + replyTo + "\r\n"
  }

  subject := msg.Subject
  if checkAscii(subject) {
    data += "Subject: " + subject + "\r\n"
  } else {
    data += "Subject: " + toUtf8(subject) + "\r\n"
  }

  if msg.IsHtml {
    data += "Content-Type: text/html; charset=UTF-8\r\n"
  } else {
    data += "Content-Type: text/plain; charset=UTF-8\r\n"
  }

  data += "\r\n"
  data += msg.Body + "\r\n"

  return data
}

func readMailMsg(mr *mail.Reader) MailMsg {
  var msg MailMsg

  // Print some info about the message
  header := mr.Header
  if msgId, err := header.MessageID(); err == nil {
    msg.MessageId = msgId
  }
  if date, err := header.Date(); err == nil {
    msg.Timestamp = date
  }
  if from, err := header.AddressList("From"); err == nil {
    msg.From = mailAddrsToString(from)
  }
  if to, err := header.AddressList("To"); err == nil {
    msg.To = mailAddrsToStrings(to)
  }
  if cc, err := header.AddressList("Cc"); err == nil {
    msg.Cc = mailAddrsToStrings(cc)
  }
  if bcc, err := header.AddressList("Bcc"); err == nil {
    msg.Bcc = mailAddrsToStrings(bcc)
  }
  if replyTo, err := header.AddressList("Reply-To"); err == nil {
    msg.ReplyTo = mailAddrsToStrings(replyTo)
  }

  if subject, err := header.Subject(); err == nil {
    msg.Subject = subject
  }

  if strings.Contains(header.Get("Content-Type"), "text/html") {
    msg.IsHtml = true
  }

  // Process each message's part
  for {
    p, err := mr.NextPart()
    if err == io.EOF {
      break
    } else if err != nil {
      panic(err)
    }

    switch h := p.Header.(type) {
    case *mail.InlineHeader:
      // This is the message's text (can be plain-text or HTML)
      b, _ := io.ReadAll(p.Body)
      msg.Body += string(b)
    case *mail.AttachmentHeader:
      // This is an attachment
      filename, _ := h.Filename()
      fmt.Printf("Got attachment: %s\n", filename)
    }
  }

  return msg
}
