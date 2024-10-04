package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadFileToS3(filePath string, bucket string, region string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file for upload: %v", err)
	}
	defer file.Close()

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

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

func StopEC2Instance(instanceID, region string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

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


func GetInstanceIDByPublicIP(publicIP string) (string, error) {
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

func StartEC2Instance(instanceID, region string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return fmt.Errorf("error creating AWS session: %v", err)
	}

	svc := ec2.New(sess)

	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	_, err = svc.StartInstances(input)
	if err != nil {
		return fmt.Errorf("error starting EC2 instance: %v", err)
	}

	log.Printf("Successfully started EC2 instance: %s", instanceID)
	return nil
}

func GetPublicIPByInstanceID(instanceID, awsRegion string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := ec2.New(sess)

	// Retry logic for fetching the public IP
	for retries := 0; retries < 24; retries++ { // Retry for up to 2 minutes (24 * 5 seconds = 120 seconds)
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(instanceID),
			},
		}

		result, err := svc.DescribeInstances(input)
		if err != nil {
			return "", fmt.Errorf("failed to describe EC2 instance: %v", err)
		}

		if len(result.Reservations) > 0 && len(result.Reservations[0].Instances) > 0 {
			instance := result.Reservations[0].Instances[0]
			if instance.PublicIpAddress != nil {
				log.Printf("Found public IP: %s for instance ID: %s", *instance.PublicIpAddress, instanceID)
				return *instance.PublicIpAddress, nil
			}
		}

		log.Printf("Public IP not yet assigned for instance ID: %s, retrying... (%d/24)", instanceID, retries+1)
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	return "", fmt.Errorf("public IP not assigned for instance ID %s after 2 minutes", instanceID)
}
