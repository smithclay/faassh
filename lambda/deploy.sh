#!/bin/zsh

# deploy.sh
# Clay Smith, 2016
# Creates deployment archive for a node.js lambda function and uploads using the AWS CLI

FUNCTION_NAME=SshFunction
DEPLOY_ARCHIVE=build/deploy.zip

echo "Creating ZIP archive..."
mkdir -p build
zip -r $DEPLOY_ARCHIVE index.js faassh bin id_rsa id_rsa.pub

echo "Updating AWS Lambda Code..."
AWS_DEFAULT_REGION=us-west-2 aws lambda update-function-code --function-name $FUNCTION_NAME --zip-file fileb://$DEPLOY_ARCHIVE
