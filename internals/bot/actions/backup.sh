#!/bin/bash

BACKUP_FILE=$1
BACKUP_PATH=$2
S3_BUCKET=$3
AWS_SECRET_ACCESS_KEY=$4
AWS_ACCESS_KEY_ID=$5
AWS_REGION=$6

if ! command -v aws &> /dev/null; then
    echo "Installing AWS CLI..."
    sudo yum install -y awscli
fi

export AWS_ACCESS_KEY_ID=$4
export AWS_SECRET_ACCESS_KEY=$5
export AWS_DEFAULT_REGION=$6


echo "Creating backup..."
tar -czf $BACKUP_FILE -C $BACKUP_PATH .

echo "Upload to s3..."
aws s3 cp $BACKUP_FILE s3://$S3_BUCKET/backup.tar.gz --region $AWS_REGION

echo "Cleanup the backup zip file"
rm -f $BACKUP_FILE