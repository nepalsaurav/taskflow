package pkg

import (
	"os"
	"os/exec"
)

var (
	POSTFIX_LOG_FILE = "/var/log/mail.log"
)

func getLogReportText() (string, error) {
	args := []string{
		POSTFIX_LOG_FILE,
		"--iso_date_time",
		"-q",
		"-smtpd-warning-detail=0",
		"--rej_add_from",
		"--verbose_msg_detail",
	}

	cmd := exec.Command("pflogsumm", args...)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), err
}

func reportToMarkdown(report string) (string, error) {

	inputfile, err := os.CreateTemp("", "pflogsumm_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(inputfile.Name())

	_, err = inputfile.Write([]byte(report))
	if err != nil {
		return "", err
	}
	inputfile.Close()

	outputfile, err := os.CreateTemp("", "pflogsumm_*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(outputfile.Name())

	pandocCommand := exec.Command("pandoc", inputfile.Name(), "-t", "markdown", "-o", outputfile.Name())

	if err := pandocCommand.Run(); err != nil {
		return "", err
	}

	md, err := os.ReadFile(outputfile.Name())

	if err != nil {
		return "", err
	}
	return string(md), nil
}

func PostfixLogDetail() (string, error) {
	report, err := getLogReportText()
	if err != nil {
		return "", err
	}
	return reportToMarkdown(report)
}
