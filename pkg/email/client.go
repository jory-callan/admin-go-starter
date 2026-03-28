package email

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"sync"
	"time"
)

// Client 邮件客户端
type Client struct {
	cfg  Config
	auth smtp.Auth
	mu   sync.RWMutex
}

// NewClient 创建邮件客户端
func NewClient(cfg Config) (*Client, error) {
	client := &Client{
		cfg: cfg,
	}

	// 初始化认证
	if cfg.Username != "" && cfg.Password != "" {
		client.auth = smtp.PlainAuth(
			"",
			cfg.Username,
			cfg.Password,
			cfg.SMTPHost,
		)
	}

	return client, nil
}

// EmailMessage 邮件消息结构
type EmailMessage struct {
	To          []string           // 收件人列表
	Cc          []string           // 抄送列表
	Bcc         []string           // 密送列表
	Subject     string             // 邮件主题
	Body        string             // 邮件正文
	ContentType string             // 内容类型：text/plain, text/html
	Attachments []*EmailAttachment // 附件列表
	Headers     map[string]string  // 自定义头
}

// EmailAttachment 邮件附件
type EmailAttachment struct {
	Filename string // 文件名
	Data     []byte // 文件内容
}

// Send 发送邮件
func (c *Client) Send(msg *EmailMessage) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 构建邮件内容
	buffer := bytes.NewBuffer(nil)

	// From header
	if c.cfg.FromName != "" {
		buffer.WriteString(fmt.Sprintf("From: \"%s\" <%s>\r\n", c.cfg.FromName, c.cfg.FromEmail))
	} else {
		buffer.WriteString(fmt.Sprintf("From: %s\r\n", c.cfg.FromEmail))
	}

	// To header
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", joinStrings(msg.To)))

	// Cc header
	if len(msg.Cc) > 0 {
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", joinStrings(msg.Cc)))
	}

	// Subject header
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))

	// Content-Type header
	contentType := msg.ContentType
	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}
	buffer.WriteString(fmt.Sprintf("Content-Type: %s\r\n\r\n", contentType))

	// Body
	buffer.WriteString(msg.Body)

	// 添加自定义 headers
	for key, value := range msg.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// 构建收件人列表（包含 To, Cc, Bcc）
	recipients := make([]string, 0)
	recipients = append(recipients, msg.To...)
	recipients = append(recipients, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", c.cfg.SMTPHost, c.cfg.SMTPPort)

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer client.Close()

	// 设置超时（通过 context 控制）
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.Timeout)*time.Second)
	defer cancel()
	_ = ctx // 用于后续扩展

	if c.auth != nil {
		if err := client.Auth(c.auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	// 设置发件人
	from := c.cfg.FromEmail
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("set from: %w", err)
	}

	// 设置收件人
	for _, to := range recipients {
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("set recipient: %w", err)
		}
	}

	// 写入邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("get data writer: %w", err)
	}
	defer w.Close()

	if _, err := w.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("write email body: %w", err)
	}

	return nil
}

// SendSimple 发送简单邮件（快捷方法）
func (c *Client) SendSimple(to []string, subject, body string) error {
	return c.Send(&EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

// SendHTML 发送 HTML 邮件
func (c *Client) SendHTML(to []string, subject, htmlBody string) error {
	return c.Send(&EmailMessage{
		To:          to,
		Subject:     subject,
		Body:        htmlBody,
		ContentType: "text/html; charset=utf-8",
	})
}

// HealthCheck 检查 SMTP 连接健康状态
func (c *Client) HealthCheck() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	addr := fmt.Sprintf("%s:%d", c.cfg.SMTPHost, c.cfg.SMTPPort)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer client.Close()

	if err := client.Noop(); err != nil {
		return fmt.Errorf("smtp noop: %w", err)
	}

	return nil
}

// Close 关闭连接（对于 SMTP 客户端通常不需要）
func (c *Client) Close() error {
	return nil
}

// joinStrings 拼接字符串数组
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += ", " + strs[i]
	}
	return result
}

func Shutdown() {

}
