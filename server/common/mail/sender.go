package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"server/api/dao"
	"strings"
	"log"
)

// SendMailAsync 异步发送邮件封装（不阻塞主进程）
func SendMailAsync(to []string, subject string, body string) {
	go func() {
		err := SendMail(to, subject, body)
		if err != nil {
			log.Printf("[Mail] Failed to send email asynchronously to %v: %v\n", to, err)
		}
	}()
}

// SendMail 发送邮件的封装
// 从系统的 SysSetting 中读取 SMTP 动态配置
func SendMail(to []string, subject string, body string) error {
	// 判断系统是否启用了 SMTP 发信功能
	enable := dao.GetSysSettingBoolValue("notification.smtp_enable", false)
	if !enable {
		return fmt.Errorf("SMTP邮件服务未启用")
	}

	hostItem, err := dao.GetSysSettingByKey("notification.smtp_host")
	if err != nil || hostItem.SettingValue == "" {
		return fmt.Errorf("SMTP主机未配置")
	}
	host := hostItem.SettingValue

	port := dao.GetSysSettingIntValue("notification.smtp_port", 465)

	userItem, err := dao.GetSysSettingByKey("notification.smtp_user")
	if err != nil || userItem.SettingValue == "" {
		return fmt.Errorf("SMTP账号未配置")
	}
	user := userItem.SettingValue

	passItem, err := dao.GetSysSettingByKey("notification.smtp_pass")
	if err != nil {
		return fmt.Errorf("SMTP密码未配置")
	}
	pass := passItem.SettingValue

	auth := smtp.PlainAuth("", user, pass, host)

	header := make(map[string]string)
	header["From"] = user
	header["To"] = strings.Join(to, ";")
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", host, port)

	// 企业级应用通常使用 465 SSL, 或者 587 STARTTLS
	if port == 465 {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return err
		}
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return err
		}
		if err = client.Mail(user); err != nil {
			return err
		}
		for _, addr := range to {
			if err = client.Rcpt(addr); err != nil {
				return err
			}
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		return client.Quit()
	} else {
		// 普通明文或 STARTTLS
		return smtp.SendMail(addr, auth, user, to, []byte(message))
	}
}
