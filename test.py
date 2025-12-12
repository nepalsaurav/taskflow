import smtplib
from email.message import EmailMessage

# Create the email
msg = EmailMessage()
msg["From"] = "nepalsaurav123@gmail.com"
msg["To"] = "nepalsaurav123@gmail.com"
msg["Subject"] = "Test Email"
msg.set_content("Hello! This is a test email sent via Postfix SMTP.")

# Connect to local Postfix SMTP server
with smtplib.SMTP("localhost", 25) as smtp:
    smtp.send_message(msg)

print("Email sent successfully!")
