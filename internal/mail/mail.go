package mail

import "github.com/wneessen/go-mail"

type Service struct {
	c    *mail.Client
	from string
}

func NewMailService(server, from, username, password string) (*Service, error) {
	c, err := mail.NewClient(server, mail.WithTLSPortPolicy(mail.TLSMandatory), mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		return nil, err
	}

	return &Service{c: c, from: from}, nil
}

func (s *Service) Send(to, subject, body string) error {
	message := mail.NewMsg()
	if err := message.From(s.from); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	if err := s.c.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
