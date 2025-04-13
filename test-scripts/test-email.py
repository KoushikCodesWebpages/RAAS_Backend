import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

# Replace with your Gmail credentials
email_user = 'koushikaltacc@gmail.com'
email_password = 'zipz kxhm kptb bqzy'  # Use your app-specific password

# Recipient email
email_to = 'koushiksample@gmail.com'

# Create the email content
subject = 'Test Email'
body = 'This is a test email sent from Python using an app-specific password!'

msg = MIMEMultipart()
msg['From'] = email_user
msg['To'] = email_to
msg['Subject'] = subject
msg.attach(MIMEText(body, 'plain'))

# Connect to Gmail's SMTP server
server = smtplib.SMTP('smtp.gmail.com', 587)
server.starttls()  # Start TLS encryption

# Login using your Gmail credentials (email and app-specific password)
server.login(email_user, email_password)

# Send the email
server.sendmail(email_user, email_to, msg.as_string())

# Quit the server connection
server.quit()

print('Email sent successfully!')
