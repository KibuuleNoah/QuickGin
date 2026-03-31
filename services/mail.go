package services

import (
	"fmt"
	"net/smtp"
	"strings"
)

// MailConfig holds SMTP connection settings loaded from environment variables.
type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// MailService sends transactional emails via SMTP.
type MailService struct {
	cfg MailConfig
}

func NewMailService(cfg MailConfig) *MailService {
	return &MailService{cfg: cfg}
}

// SendOTP sends a 6-digit OTP to the given e-mail address.
func (m *MailService) SendOTP(to, otp string) error {
	subject := "Your verification code"
	body := buildOTPEmail(otp)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", m.cfg.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

func buildOTPEmail(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:32px">
  <h2 style="margin-bottom:8px">Verification code</h2>
  <p style="color:#555">Use the code below to sign in. It expires in 5 minutes.</p>
  <div style="font-size:36px;font-weight:700;letter-spacing:8px;margin:24px 0;color:#111">%s</div>
  <p style="color:#999;font-size:13px">If you didn't request this, you can safely ignore this email.</p>
</body>
</html>`, otp)
}
