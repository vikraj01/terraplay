package game

import (
	"fmt"
	"log"
	"os"

	"github.com/vikraj01/terraplay/internals/utils"
)

func BackupAndStopEC2(sshConfig utils.SSHConfig, backupPath string, s3Bucket string, backupFile string, publicIP string) error {
	client, err := utils.ConnectToEC2ViaSSH(sshConfig)
	if err != nil {
		return fmt.Errorf("error connecting to EC2 via SSH: %v", err)
	}
	defer client.Close()

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
	err = utils.RunCommandOnEC2(client, createScriptCommand)
	if err != nil {
		return fmt.Errorf("error creating backup.sh script on EC2: %v", err)
	}
	log.Println("Backup script created successfully on EC2 instance.")

	backupCommand := fmt.Sprintf("/bin/bash /home/ec2-user/backup.sh %s %s %s %s %s %s", backupFile, backupPath, s3Bucket, awsSecretKey, awsAccessKey, awsRegion)
	err = utils.RunCommandOnEC2(client, backupCommand)
	if err != nil {
		return fmt.Errorf("error running backup script on EC2: %v", err)
	}
	log.Println("Backup and upload completed successfully on EC2 instance.")

	instanceID, err := utils.GetInstanceIDByPublicIP(publicIP)
	if err != nil {
		return fmt.Errorf("error retrieving instance ID: %v", err)
	}

	if err := utils.StopEC2Instance(instanceID, awsRegion); err != nil {
		return fmt.Errorf("error stopping EC2 instance: %v", err)
	}

	return nil
}

// halted, terminated, running, pending