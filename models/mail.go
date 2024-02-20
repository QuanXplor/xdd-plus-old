package models

import (
	"errors"
	"net/smtp"
	"strconv"
	"strings"
)

type Mail struct {
	Server   string
	Port     int
	Sender   string
	Password string
	Name     string
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func SendToMail(mail Mail, to, subject, body string) error {
	if mail.Name == "" {
		mail.Name = mail.Sender
	}
	mailType := ""
	auth := LoginAuth(mail.Sender, mail.Password)
	var content_type string
	if mailType == "html" {
		content_type = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + mail.Name + "<" + mail.Server + ">" + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(mail.Server+":"+strconv.Itoa(mail.Port), auth, mail.Sender, send_to, msg)
	return err
}
