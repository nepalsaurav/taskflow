package pkg

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
)

type PostfixConfig struct {
	SMTPAccount struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Auth     string `json:"auth"`
		User     string `json:"user"`
		Password string `json:"password"`
		From     string `json:"from"`
	} `json:"smtpAccount"`
}

func SetPostfixConfig(postfixConfig PostfixConfig) (PostfixConfig, error) {
	tmpl, err := template.ParseFiles("conf/postfix.gotmpl")
	if err != nil {
		log.Printf("error parsing postfix config file, err: %v\n", err)
		return PostfixConfig{}, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("error on getting hostname, err: %v\n", err)
		return PostfixConfig{}, err
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Printf("error on getting current user, err: %v\n", err)
		return PostfixConfig{}, err
	}
	data := map[string]string{
		"Hostname":      hostname,
		"RelayHost":     postfixConfig.SMTPAccount.Host,
		"RelayHostPort": strconv.Itoa(postfixConfig.SMTPAccount.Port),
		"HostUserName":  currentUser.Name,
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, data)
	if err != nil {
		log.Printf("error on executing template with data, err: %v\n", err)
		return PostfixConfig{}, err
	}

	if err := addSMTPPassword(postfixConfig); err != nil {
		return PostfixConfig{}, err
	}

	// write posfixconfig using sudo + tee
	cmd := exec.Command("sudo", "tee", postfixConfigPath)
	cmd.Stdin = bytes.NewReader(buff.Bytes())

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to write postfix config in %s error: %v\nOutput: %s", postfixConfigPath, err, output)
		return PostfixConfig{}, err
	}

	// reload postfix
	if out, err := exec.Command("sudo", "postfix", "reload").CombinedOutput(); err != nil {
		log.Printf("reload posfix failed err: %v\n%s", err, out)
		return PostfixConfig{}, err
	}

	return postfixConfig, nil
}

func addSMTPPassword(acc PostfixConfig) error {
	entry := fmt.Sprintf("[%s]:%d\t%s:%s\n", acc.SMTPAccount.Host, acc.SMTPAccount.Port, acc.SMTPAccount.User, acc.SMTPAccount.Password)

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
