#!/bin/bash
yum update -y

yum install -y httpd git curl

systemctl start httpd
systemctl enable httpd

echo "<html><body><h1>Hello from Amazon Linux 2!</h1></body></html>" > /var/www/html/index.html

chown -R apache:apache /var/www/html

firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --reload
