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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
		User:       "ec2-user",                              // Change as necessary
		PrivateKey: privateKey, // Load from environment or securely
	}

	backupPath := "/opt/minetest/data" // Adjust the path you need to backup
	s3Bucket := "global-bucket-893606"

	err = BackupAndStopEC2(sshConfig, backupPath, s3Bucket, details.ServerIP)
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

func BackupAndStopEC2(sshConfig SSHConfig, backupPath string, s3Bucket string, publicIP string) error {
	client, err := connectToEC2ViaSSH(sshConfig)
	if err != nil {
		return fmt.Errorf("error connecting to EC2 via SSH: %v", err)
	}
	defer client.Close()

	backupFile := "/tmp/backup.tar.gz"
	command := fmt.Sprintf("tar -czf %s -C %s .", backupFile, backupPath)
	if err := runCommandOnEC2(client, command); err != nil {
		return fmt.Errorf("error creating backup archive: %v", err)
	}

	if err := uploadFileToS3(backupFile, s3Bucket); err != nil {
		return fmt.Errorf("error uploading backup to S3: %v", err)
	}

	instanceID, err := getInstanceIDByPublicIP(publicIP)
	if err != nil {
		return fmt.Errorf("error retrieving instance ID: %v", err)
	}

	if err := stopEC2Instance(instanceID); err != nil {
		return fmt.Errorf("error stopping EC2 instance: %v", err)
	}

	return nil
}

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
