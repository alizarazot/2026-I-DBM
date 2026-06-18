package mail

import (
	"bytes"
	"strings"

	_ "embed"

	"github.com/signintech/gopdf"
	"github.com/wneessen/go-mail"
)

//go:embed font.ttf
var font []byte

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
	message, err := s.makeMessage(to, subject, body)
	if err != nil {
		return err
	}
	if err := s.c.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (s *Service) SendPDF(to, subject, body, pdf string) error {
	message, err := s.makeMessage(to, subject, body)
	if err != nil {
		return err
	}

	p := gopdf.GoPdf{}
	p.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err = p.AddTTFFontData("liberation", font)
	if err != nil {
		return err
	}

	err = p.SetFont("liberation", "", 14)
	if err != nil {
		return err
	}

	p.AddPage()

	texts := strings.Split(pdf, "\n\n")
	for i, t := range texts {
		if i%2 == 0 {
		}
		p.SetXY(50, float64(50*(i+1)))
		p.Text(t)
	}

	var buf bytes.Buffer
	if _, err := p.WriteTo(&buf); err != nil {
		return err
	}

	if err := message.AttachReader("doc.pdf", &buf); err != nil {
		return err
	}

	if err := s.c.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (s *Service) makeMessage(to, subject, body string) (*mail.Msg, error) {
	message := mail.NewMsg()
	if err := message.From(s.from); err != nil {
		return nil, err
	}
	if err := message.To(to); err != nil {
		return nil, err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	return message, nil
}
