package emailsend

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/sansanbaby/dayreport/config"
)

type EmailSender struct {
	config config.EmailConfig
}

// 创建一个邮件发送器
func NewEmailSender() *EmailSender {
	return &EmailSender{
		config: config.Email,
	}
}

// 发送纯邮件
//
//	func (e *EmailSender) SendSubjectAndBody(subject, body string) error {
//		header := make(map[string]string)
//		header["From"] = e.config.From
//		header["To"] = strings.Join(e.config.To, ",")
//		header["Subject"] = subject
//		header["MIME-Version"] = "1.0"
//		header["Content-Type"] = "text/html; charset=UTF-8"
//
//		message := ""
//		for k, v := range header {
//			message += fmt.Sprintf("%s: %s\r\n", k, v)
//		}
//		message += "\r\n" + body
//
//		auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.SMTPServer)
//		addr := fmt.Sprintf("%s:%d", e.config.SMTPServer, e.config.SMTPPort)
//
//		tlsConfig := &tls.Config{
//			InsecureSkipVerify: true,
//			ServerName:         e.config.SMTPServer,
//		}
//
//		conn, err := tls.Dial("tcp", addr, tlsConfig)
//		if err != nil {
//			return fmt.Errorf("连接 SMTP 服务器失败：%v", err)
//		}
//		defer conn.Close()
//
//		client, err := smtp.NewClient(conn, e.config.SMTPServer)
//		if err != nil {
//			return fmt.Errorf("创建 SMTP 客户端失败：%v", err)
//		}
//		defer client.Close()
//
//		if err = client.Auth(auth); err != nil {
//			return fmt.Errorf("SMTP 认证失败：%v", err)
//		}
//
//		if err = client.Mail(e.config.From); err != nil {
//			return fmt.Errorf("设置发件人失败：%v", err)
//		}
//
//		for _, to := range e.config.To {
//			if err = client.Rcpt(to); err != nil {
//				return fmt.Errorf("设置收件人失败：%v", err)
//			}
//		}
//
//		w, err := client.Data()
//		if err != nil {
//			return fmt.Errorf("获取数据写入器失败：%v", err)
//		}
//
//		_, err = w.Write([]byte(message))
//		if err != nil {
//			return fmt.Errorf("写入邮件内容失败：%v", err)
//		}
//
//		err = w.Close()
//		if err != nil {
//			return fmt.Errorf("关闭数据写入器失败：%v", err)
//		}
//
//		return client.Quit()
//	}
//
// 发送带附件的邮件

// 发送带附件的邮件
func (e *EmailSender) SendWithAttachment(subject, body, attachmentPath string) error {
	fileData, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("读取附件失败：%v", err)
	}

	encodedFile := base64.StdEncoding.EncodeToString(fileData)
	filename := filepath.Base(attachmentPath)

	message := fmt.Sprintf("--BOUNDARY\r\n")
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "Content-Transfer-Encoding: quoted-printable\r\n\r\n"
	message += body + "\r\n"
	message += fmt.Sprintf("\r\n--BOUNDARY\r\n")
	message += "Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet; name=\"" + filename + "\"\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"
	message += "Content-Disposition: attachment; filename=\"" + filename + "\"\r\n\r\n"

	lineLength := 76
	for i := 0; i < len(encodedFile); i += lineLength {
		end := i + lineLength
		if end > len(encodedFile) {
			end = len(encodedFile)
		}
		message += encodedFile[i:end] + "\r\n"
	}

	message += "\r\n--BOUNDARY--"

	headers := make(map[string]string)
	headers["From"] = e.config.From
	headers["To"] = strings.Join(e.config.To, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=BOUNDARY"

	finalMessage := ""
	for k, v := range headers {
		finalMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	finalMessage += "\r\n" + message

	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.SMTPServer)
	addr := fmt.Sprintf("%s:%d", e.config.SMTPServer, e.config.SMTPPort)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.config.SMTPServer,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接 SMTP 服务器失败：%v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.config.SMTPServer)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败：%v", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败：%v", err)
	}

	if err = client.Mail(e.config.From); err != nil {
		return fmt.Errorf("设置发件人失败：%v", err)
	}

	for _, to := range e.config.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("设置收件人失败：%v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败：%v", err)
	}

	_, err = w.Write([]byte(finalMessage))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败：%v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入器失败：%v", err)
	}

	return client.Quit()
}
