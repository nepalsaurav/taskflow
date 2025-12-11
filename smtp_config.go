package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const SMTP_CONFIG_FILE = "config/smtp_config.json"

type SMTPAccount struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Auth     string `json:"auth"`
	User     string `json:"user"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type SMTPConfig struct {
	Accounts []SMTPAccount `json:"accounts"`
}

// readSMTPConfigFile reads the JSON config and returns SMTPConfig
func readSMTPConfigFile() SMTPConfig {
	data, err := os.ReadFile(SMTP_CONFIG_FILE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", SMTP_CONFIG_FILE, err)
		os.Exit(1)
	}

	var cfg SMTPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "JSON parse error: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

// checkPostfixRunning ensures Postfix master process is running
func checkPostfixRunning() {
	cmd := exec.Command("ps", "-C", "master", "-o", "pid=")
	out, _ := cmd.Output()
	if strings.TrimSpace(string(out)) == "" {
		fmt.Fprintf(os.Stderr, "Postfix is not running, please install or start postfix\n")
		os.Exit(1)
	}
}

// runPostConf updates a Postfix parameter using postconf
func runPostConf(key, value string) {
	checkPostfixRunning()
	cmd := exec.Command("sudo", "postconf", "-e", fmt.Sprintf("%s=%s", key, value))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "postconf error for %s: %v\n", key, err)
		os.Exit(1)
	}
}

func main() {
	cfg := readSMTPConfigFile()
	if len(cfg.Accounts) == 0 {
		fmt.Fprintln(os.Stderr, "No SMTP accounts found")
		os.Exit(1)
	}
	acc := cfg.Accounts[0]

	// Create Postfix SASL directory
	homeDir, _ := os.UserHomeDir()
	passFile := filepath.Join(homeDir, "postfix", "sasl_passwd")
	if err := os.MkdirAll(filepath.Dir(passFile), 0700); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create directory: %v\n", err)
		os.Exit(1)
	}

	// Write plain SASL credentials
	passContent := fmt.Sprintf("[%s]:%d %s:%s\n", acc.Host, acc.Port, acc.User, acc.Password)
	if err := os.WriteFile(passFile, []byte(passContent), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "cannot write sasl_password: %v\n", err)
		os.Exit(1)
	}

	// Hash the password map for Postfix
	cmd := exec.Command("postmap", passFile)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "postmap error: %v\n", err)
		os.Exit(1)
	}

	// Configure Postfix relay
	// configure postfix relay using postconf
	runPostConf("relayhost", fmt.Sprintf("[%s]:%d", acc.Host, acc.Port))
	runPostConf("smtp_sasl_auth_enable", "yes")
	runPostConf("smtp_sasl_password_maps", passFile)
	runPostConf("smtp_sasl_security_options", "noanonymous")
	runPostConf("smtp_tls_security_level", "encrypt")

	// Reload postfix to apply changes
	reload := exec.Command("sudo", "postfix", "reload")
	if err := reload.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "postfix reload error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Postfix relay configuration applied successfully")
}
