package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

const (
	postfixConfigPath = "/etc/postfix/main.cf"
	passwdFile        = "/etc/postfix/sasl/sasl_passwd"
	dbFile            = "/etc/postfix/sasl/sasl_passwd.db"
)

type SMTPAccount struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Auth     string `json:"auth"`
	User     string `json:"user"`
	Password string `json:"password"`
	From     string `json:"from"`
}

func SetPostfixConfig(smtpAccount SMTPAccount) error {
	tmpl, err := template.ParseFiles("conf/postfix.gotmpl")
	if err != nil {
		log.Printf("error parsing postfix config file, err: %v\n", err)
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("error on getting hostname, err: %v\n", err)
		return err
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Printf("error on getting current user, err: %v\n", err)
		return err
	}
	data := map[string]string{
		"Hostname":      hostname,
		"RelayHost":     smtpAccount.Host,
		"RelayHostPort": strconv.Itoa(smtpAccount.Port),
		"HostUserName":  currentUser.Name,
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, data)
	if err != nil {
		log.Printf("error on executing template with data, err: %v\n", err)
		return err
	}

	if err := addSMTPPassword(smtpAccount); err != nil {
		return err
	}

	// write posfixconfig using sudo + tee
	cmd := exec.Command("sudo", "tee", postfixConfigPath)
	cmd.Stdin = bytes.NewReader(buff.Bytes())

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to write postfix config in %s error: %v\nOutput: %s", postfixConfigPath, err, output)
		return err
	}

	// reload postfix
	if out, err := exec.Command("sudo", "postfix", "reload").CombinedOutput(); err != nil {
		log.Printf("reload posfix failed err: %v\n%s", err, out)
		return err
	}

	return nil
}

func addSMTPPassword(acc SMTPAccount) error {
	entry := fmt.Sprintf("[%s]:%d\t%s:%s\n", acc.Host, acc.Port, acc.User, acc.Password)

	// 2. Write to /etc/postfix/sasl/sasl_passwd
	cmd := exec.Command("sudo", "tee", passwdFile)
	cmd.Stdin = strings.NewReader(entry)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to write sasl_passwd: %v\n%s", err, out)
	}

	// 2. Fix permissions, rebuild DB
	commands := [][]string{
		{"sudo", "chmod", "600", "/etc/postfix/sasl/sasl_passwd"},
		{"sudo", "postmap", "/etc/postfix/sasl/sasl_passwd"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed setting smtp password %v\n%s", err, out)

		}
	}
	return nil
}
