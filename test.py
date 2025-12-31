import smtplib
from concurrent.futures import ThreadPoolExecutor
from email.message import EmailMessage


# Email sending function
def send_email(index):
    sender = "nepalsaurav123@gmail.com"
    recipient = "nepalsaurav123@gmail.com,nepalsaurav@yandex.com"
    subject = f"Test Email {index} from Postfix"
    body = f"Hello! This is test email number {index} sent through local Postfix."

    msg = EmailMessage()
    msg["From"] = sender
    msg["To"] = recipient
    msg["Subject"] = subject
    msg["X-Tracking-ID"] = "tracking_id"
    msg.set_content(body)

    try:
        with smtplib.SMTP("localhost", 25) as server:
            server.send_message(msg)
        print(f"Email {index} sent successfully!")
    except Exception as e:
        print(f"Failed to send email {index}: {e}")


# Send 20 emails concurrently
with ThreadPoolExecutor(max_workers=20) as executor:
    for i in range(1, 2):
        executor.submit(send_email, i)
