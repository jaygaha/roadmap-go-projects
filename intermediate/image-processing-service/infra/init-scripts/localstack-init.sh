#!/bin/bash

# Wait for localstack to be ready
echo "Waiting for LocalStack to be ready..."
sleep 5 

# Create S3 bucket
echo "Creating S3 bucket..."
awslocal --endpoint-url=http://localhost:4566 s3 mb s3://image-processing-bucket

# Set bucket policy for public read access (optional)
echo "Setting bucket policy for public read access..."
awslocal --endpoint-url=http://localhost:4566 s3api put-bucket-policy --bucket image-processing-bucket --policy '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::image-processing-bucket/*"
    }
  ]
}'

# List buckets to confirm creation
echo "Listing S3 buckets:"
awslocal s3 ls

echo "LocalStack initialization completed!"