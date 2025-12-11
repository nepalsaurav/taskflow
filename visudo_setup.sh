#!/bin/bash
set -e
USERNAME="$(whoami)"
SUDOERS_FILE="/etc/sudoers.d/postfix-taskflow"

echo "Updating sudoers rule for full Postfix + SMTP password support..."

sudo tee "$SUDOERS_FILE" > /dev/null <<EOF
# Postfix config management for $USERNAME (taskflow)
$USERNAME ALL=(root) NOPASSWD: /usr/bin/tee /etc/postfix/main.cf
$USERNAME ALL=(root) NOPASSWD: /usr/bin/tee /etc/postfix/main.cf.new
$USERNAME ALL=(root) NOPASSWD: /usr/bin/tee /etc/postfix/sasl/sasl_passwd
$USERNAME ALL=(root) NOPASSWD: /usr/sbin/postfix reload
$USERNAME ALL=(root) NOPASSWD: /usr/sbin/postfix check
$USERNAME ALL=(root) NOPASSWD: /usr/sbin/postmap /etc/postfix/sasl/sasl_passwd
EOF

sudo chmod 0440 "$SUDOERS_FILE"
echo "Updated sudoers rule installed: $SUDOERS_FILE"
