#!/bin/zsh

# Creates deployment archive for a node.js lambda function and uploads using the AWS CLI

mkdir -p build
zip -r build/testInVPC.zip index.js faassh id_rsa id_rsa.pub

AWS_DEFAULT_REGION=us-west-2 aws lambda update-function-code --function-name SshFunction --zip-file fileb://build/testInVPC.zip
