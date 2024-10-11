package game

import (
	"fmt"
	"log"
	"os"

	"github.com/vikraj01/terraplay/internals/utils"
)

func RestartEC2(instanceId string) (string, error) {
	err := utils.StartEC2Instance(instanceId, os.Getenv("AWS_REGION"))
	if err != nil {
		log.Printf("Error starting EC2 instance: %v", err)
		return "", err
	}

	awsRegion := os.Getenv("AWS_REGION")
	newServerIP, err := utils.GetPublicIPByInstanceID(instanceId, awsRegion)
	if err != nil {
		log.Printf("Error retrieving new server IP: %v", err)
		return "", err
	}

	return newServerIP, nil
}

func RestoreEC2(sshConfig utils.SSHConfig, backupPath, s3Bucket, backupFile, publicIP string, instanceID string) error {
	awsRegion := os.Getenv("AWS_REGION")

	err := utils.WaitForInstanceRunning(instanceID, awsRegion)
	if err != nil {
		return fmt.Errorf("error waiting for instance to reach running state: %v", err)
	}

	client, err := utils.ConnectToEC2ViaSSHWithRetry(sshConfig)
	if err != nil {
		return fmt.Errorf("error connecting to EC2 via SSH: %v", err)
	}
	defer client.Close()

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

		echo "Downloading backup from S3 bucket $S3_BUCKET..."
		aws s3 cp s3://$S3_BUCKET/backup.tar.gz $BACKUP_FILE --region $AWS_REGION

		echo "Extracting backup..."
		tar -xzf $BACKUP_FILE -C $BACKUP_PATH

		echo "Cleaning up the backup file..."
		rm -f $BACKUP_FILE
	`

	createScriptCommand := fmt.Sprintf("echo '%s' > /home/ec2-user/restore.sh && chmod +x /home/ec2-user/restore.sh", scriptContent)
	err = utils.RunCommandOnEC2(client, createScriptCommand)
	if err != nil {
		return fmt.Errorf("error creating restore.sh script on EC2: %v", err)
	}
	log.Println("Restore script created successfully on EC2 instance.")

	restoreCommand := fmt.Sprintf("/bin/bash /home/ec2-user/restore.sh %s %s %s %s %s %s", backupFile, backupPath, s3Bucket, os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_REGION"))
	err = utils.RunCommandOnEC2(client, restoreCommand)
	if err != nil {
		return fmt.Errorf("error running restore script on EC2: %v", err)
	}

	log.Println("Data restore completed successfully on EC2 instance.")
	return nil
}