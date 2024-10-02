#!/bin/bash
# Update all installed packages
yum update -y

# Install common packages
yum install -y httpd git curl

# Start and enable Apache HTTP server (httpd)
systemctl start httpd
systemctl enable httpd

# Write a simple HTML file to the web root
echo "<html><body><h1>Hello from Amazon Linux 2!</h1></body></html>" > /var/www/html/index.html

# Set permissions on the web root
chown -R apache:apache /var/www/html

# Open port 80 in the firewall (if applicable)
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --reload
