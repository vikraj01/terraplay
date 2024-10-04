package bot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bwmarrin/discordgo"
	"github.com/vikraj01/terraplay/internals/dynamodb"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host       string
	Port       string
	User       string
	PrivateKey []byte
}

func handleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	dynamodbService, err := dynamodb.InitializeDynamoDB()
	if err != nil {
		log.Printf("Error initializing DynamoDB: %v", err)
		s.ChannelMessageSend(m.ChannelID, "⚠️ Error: Could not initialize database. Please try again later.")
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Usage: `!stop <session_id>`")
		return
	}

	sessionId := args[2]
	details, err := dynamodbService.GetDetailsBySessionID(sessionId)
	if err != nil {
		log.Printf("Error fetching sessionId details: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: Could not find workspace for the given session ID: %v", err))
		return
	}
	sshKeyBase64 := os.Getenv("EC2_SSH_KEY_BASE64")
	if sshKeyBase64 == "" {
		log.Fatal("EC2_SSH_KEY_BASE64 is not set")
	}

	privateKey, err := base64.StdEncoding.DecodeString(sshKeyBase64)
	if err != nil {
		log.Fatalf("Error decoding base64 private key: %v", err)
	}
	sshConfig := SSHConfig{
		Host:       details.ServerIP,
		Port:       "22",
		User:       "ec2-user", // Change as necessary
		PrivateKey: privateKey, // Load from environment or securely
	}

	backupFile := "/tmp/backup.tar.gz"
	backupPath := "/opt/minetest/data"
	s3Bucket := "global-bucket-893606"

	err = BackupAndStopEC2(sshConfig, backupPath, s3Bucket, backupFile, details.ServerIP)
	if err != nil {
		log.Printf("Error executing backup and stop: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ Error: %v", err))
		return
	}

	message := fmt.Sprintf(
		"EC2 instance with IP `%s` has been backed up and stopped. Workspace: `%s`", details.ServerIP, details.Workspace)
	dynamodbService.UpdateSessionStatusAndIP(sessionId, "halted", details.ServerIP)
	s.ChannelMessageSend(m.ChannelID, message)
}

func connectToEC2ViaSSH(config SSHConfig) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(config.PrivateKey)

	if err != nil {
		return nil, fmt.Errorf("error parsing SSH private key: %v", err)
	}

	sshConfig := ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), &sshConfig)

	if err != nil {
		return nil, fmt.Errorf("error dailing SSH:%v", err)
	}

	return client, nil
}

func BackupAndStopEC2(sshConfig SSHConfig, backupPath string, s3Bucket string, backupFile string, publicIP string) error {
	client, err := connectToEC2ViaSSH(sshConfig)
	if err != nil {
		return fmt.Errorf("error connecting to EC2 via SSH: %v", err)
	}
	defer client.Close()

	testCommand := "echo 'SSH connection successful' > /tmp/test_ssh_connection.txt"
	err = runCommandOnEC2(client, testCommand)
	if err != nil {
		return fmt.Errorf("error running test command on EC2: %v", err)
	}
	log.Println("Test command executed successfully: /tmp/test_ssh_connection.txt created")

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	scriptContent := `
		#!/bin/bash
		BACKUP_FILE=$1
		BACKUP_PATH=$2
		S3_BUCKET=$3
		AWS_SECRET_ACCESS_KEY=$4
		AWS_ACCESS_KEY_ID=$5
		AWS_REGION=$6

		if ! command -v aws &> /dev/null; then
			echo "AWS CLI not found, installing..."
			sudo yum install -y awscli
		fi

		export AWS_ACCESS_KEY_ID=$5
		export AWS_SECRET_ACCESS_KEY=$4
		export AWS_DEFAULT_REGION=$6

		echo "Creating backup..."
		tar -czf $BACKUP_FILE -C $BACKUP_PATH .

		echo "Uploading backup to S3 bucket $S3_BUCKET..."
		aws s3 cp $BACKUP_FILE s3://$S3_BUCKET/backup.tar.gz --region $AWS_REGION

		echo "Cleaning up the backup file..."
		rm -f $BACKUP_FILE
	`

	createScriptCommand := fmt.Sprintf("echo '%s' > /home/ec2-user/backup.sh && chmod +x /home/ec2-user/backup.sh", scriptContent)
	err = runCommandOnEC2(client, createScriptCommand)
	if err != nil {
		return fmt.Errorf("error creating backup.sh script on EC2: %v", err)
	}
	log.Println("Backup script created successfully on EC2 instance.")

	backupCommand := fmt.Sprintf("/bin/bash /home/ec2-user/backup.sh %s %s %s %s %s %s", backupFile, backupPath, s3Bucket, awsSecretKey, awsAccessKey, awsRegion)
	err = runCommandOnEC2(client, backupCommand)
	if err != nil {
		return fmt.Errorf("error running backup script on EC2: %v", err)
	}
	log.Println("Backup and upload completed successfully on EC2 instance.")

	instanceID, err := getInstanceIDByPublicIP(publicIP)
	if err != nil {
		return fmt.Errorf("error retrieving instance ID: %v", err)
	}

	if err := stopEC2Instance(instanceID); err != nil {
		return fmt.Errorf("error stopping EC2 instance: %v", err)
	}

	return nil
}

func getInstanceIDByPublicIP(publicIP string) (string, error) {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("ip-address"),
				Values: []*string{aws.String(publicIP)},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil && *instance.PublicIpAddress == publicIP {
				return *instance.InstanceId, nil
			}
		}
	}

	return "", fmt.Errorf("no instance found with public IP: %s", publicIP)
}

func stopEC2Instance(instanceID string) error {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	_, err := svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		return fmt.Errorf("error stopping EC2 instance: %v", err)
	}

	log.Printf("Successfully stopped EC2 instance: %s", instanceID)
	return nil
}

// halted, terminated, running, pending

func runCommandOnEC2(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error creating SSH session: %v", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return fmt.Errorf("error running command: %v (stderr: %s)", err, stderr.String())
	}

	log.Printf("Command output: %s", stdout.String())
	return nil
}

/*

func uploadFileToS3(filePath string, bucket string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file for upload: %v", err)
	}
	defer file.Close()

	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("backup.tar.gz"),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to S3: %v", err)
	}

	log.Printf("File successfully uploaded to S3 bucket %s", bucket)
	return nil
}

*/
