package enaga

import (
  "errors"
  "fmt"
  "github.com/emersion/go-imap"
  "github.com/emersion/go-imap/backend"
  "github.com/emersion/go-imap/server"
  "github.com/emersion/go-smtp"
  "io"
  "log"
  "time"
)

// The Backend implements SMTP server methods.
type Backend struct {
  users map[string]*backend.User
}

func (bkd *Backend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {

  return nil, errors.New("Bad username or password")
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(conn *smtp.Conn) (smtp.Session, error) {
  return &Session{}, nil
}

// A Session is returned after successful login.
type Session struct{}

func (s Session) Reset() {
}

func (s Session) Logout() error {
  return nil
}

func (s Session) AuthPlain(username, password string) error {

  return nil
}

func (s Session) Mail(from string, opts *smtp.MailOptions) error {

  return nil
}

func (s Session) Rcpt(to string, opts *smtp.RcptOptions) error {

  return nil
}

func (s Session) Data(r io.Reader) error {
  if b, err := io.ReadAll(r); err != nil {
    return err
  } else {
    fmt.Println(string(b))
  }
  return nil
}

func SmtpListener(domain string, port int16, timeout int64) {
  be := &Backend{}

  s := smtp.NewServer(be)
  s.Addr = fmt.Sprintf("%s:%d", domain, port)
  s.Domain = domain
  s.WriteTimeout = time.Duration(timeout) * time.Second
  s.ReadTimeout = time.Duration(timeout) * time.Second
  s.MaxMessageBytes = 5 * 1024 * 1024
  s.MaxRecipients = 5000
  s.AllowInsecureAuth = false

  fmt.Println("Starting server at", s.Addr)
  if err := s.ListenAndServe(); err != nil {
    panic(err)
  }

}

func ImapListener(domain string, port int16) {
  be := &Backend{}

  // Create a new server
  s := server.New(be)
  s.Addr = fmt.Sprintf("%s:%d", domain, port)
  // Since we will use this server for testing only, we can allow plain text
  // authentication over unencrypted connections
  s.AllowInsecureAuth = false

  log.Println("Starting IMAP server at ", fmt.Sprintf("%s:%d", domain, port))
  if err := s.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
