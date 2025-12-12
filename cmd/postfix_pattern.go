// cmd/postfix_patterns.go
package cmd

var PostfixPatternDefinitions = map[string]string{
	// ────────────────────────────── Helpers ──────────────────────────────
	"GREEDYDATA_NO_COLON":     `[^:]*`,
	"GREEDYDATA_NO_SEMICOLON": `[^;]*`,
	"GREEDYDATA_NO_BRACKET":   `[^<>]*`,
	"STATUS_WORD":             `[\\w-]*`,
	"IP_UNKNOWN":              `unknown`,

	// ────────────────────────────── Core ──────────────────────────────
	"POSTFIX_QUEUEID":       `([0-9A-F]{6,}|[0-9a-zA-Z]{12,}|NOQUEUE)`,
	"POSTFIX_CLIENT":        `%{HOSTNAME:postfix_client_hostname}?\\[%{IP_UNKNOWN:postfix_client_ip_unknown}|%{IP:postfix_client_ip}\\](?::%{INT:postfix_client_port})?`,
	"POSTFIX_RELAY":         `%{HOSTNAME:postfix_relay_hostname}?\\[%{IP:postfix_relay_ip}|%{DATA:postfix_relay_service}\\](?::%{INT:postfix_relay_port})?|%{WORD:postfix_relay_service}`,
	"POSTFIX_SMTP_STAGE":    `(CONNECT|HELO|EHLO|STARTTLS|AUTH|MAIL( FROM)?|RCPT( TO)?|(end of )?DATA|BDAT|RSET|UNKNOWN|END-OF-MESSAGE|VRFY|\\.)`,
	"POSTFIX_ACTION":        `(accept|defer|discard|filter|header-redirect|milter-reject|reject|reject_warning)`,
	"POSTFIX_KEYVALUE_DATA": `[\\w-]+=[^;]*`,
	"POSTFIX_KEYVALUE":      `%{POSTFIX_QUEUEID:postfix_queueid}: %{POSTFIX_KEYVALUE_DATA:postfix_keyvalue_data}`,

	// TLS connection line (RE2-safe version – no (?<name>) and no complex nested groups)
	"POSTFIX_TLSCONN": `%{DATA:postfix_tls_level} TLS connection established (?:to %{POSTFIX_RELAY:postfix_tls_peer}|from %{POSTFIX_CLIENT:postfix_tls_peer}): %{DATA:postfix_tls_version} with cipher %{DATA:postfix_tls_cipher} \\(%{INT:postfix_tls_bits} bits\\)`,

	// Warning lines (RE2-safe – removed illegal (?<name>…))
	"POSTFIX_WARNING": `(%{POSTFIX_QUEUEID:postfix_queueid}: )?(warning|fatal|info): [^;]+`,

	// ────────────────────────────── Message types ──────────────────────────────
	"POSTFIX_SMTPD_CONNECT":      `connect from %{POSTFIX_CLIENT:postfix_client}`,
	"POSTFIX_SMTPD_DISCONNECT":   `disconnect from %{POSTFIX_CLIENT:postfix_client}`,
	"POSTFIX_SMTPD_LOSTCONN":     `%{DATA:postfix_lost_reason} from %{POSTFIX_CLIENT:postfix_client}`,
	"POSTFIX_SMTPD_NOQUEUE":      `%{POSTFIX_QUEUEID:postfix_queueid}: %{POSTFIX_ACTION:postfix_action}: %{POSTFIX_SMTP_STAGE:postfix_stage} from %{POSTFIX_CLIENT:postfix_client}`,
	"POSTFIX_CLEANUP_MESSAGEID":  `%{POSTFIX_QUEUEID:postfix_queueid}: message-id=<?%{GREEDYDATA_NO_BRACKET:postfix_message_id}>?`,
	"POSTFIX_QMGR_REMOVED":       `%{POSTFIX_QUEUEID:postfix_queueid}: removed`,
	"POSTFIX_SMTP_DELIVERY":      `%{POSTFIX_KEYVALUE} status=%{STATUS_WORD:postfix_status}( \\(%{GREEDYDATA:postfix_response}\\))?`,
	"POSTFIX_POSTSCREEN_CONNECT": `CONNECT from %{POSTFIX_CLIENT:postfix_client} to \\[%{IP:postfix_server_ip}\\]:%{INT:postfix_server_port}`,

	// ──────────────────────── THIS IS THE ONLY ONE YOU COMPILE ────────────────────────
	"POSTFIX_LINE": `^%{SYSLOGBASE} (?:
          %{POSTFIX_SMTPD_CONNECT}
        | %{POSTFIX_SMTPD_DISCONNECT}
        | %{POSTFIX_SMTPD_LOSTCONN}
        | %{POSTFIX_SMTPD_NOQUEUE}
        | %{POSTFIX_CLEANUP_MESSAGEID}
        | %{POSTFIX_QMGR_REMOVED}
        | %{POSTFIX_SMTP_DELIVERY}
        | %{POSTFIX_POSTSCREEN_CONNECT}
        | %{POSTFIX_TLSCONN}
        | %{POSTFIX_WARNING}
        | %{POSTFIX_QUEUEID:postfix_queueid}:.*
        | %{GREEDYDATA:postfix_message}
    )$`,
}
