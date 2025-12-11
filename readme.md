# Taskflow Postfix Configuration

> **Note:** This project is in an **early stage** and intended for experimental use.  

This Go package automates configuring Postfix with a single SMTP relay server. It generates `main.cf`, manages SASL passwords, and reloads Postfix.  

---

## Requirements

- Ubuntu/Debian system
- Postfix installed (local-only setup is recommended)
- Go 1.20+ (or compatible)
- Root or sudo access to manage Postfix configuration

---

## Setup Steps

1. **Install Postfix locally**

```bash
sudo apt update
sudo apt install postfix
