#!/bin/zsh

aws lambda invoke --invocation-type RequestResponse --function-name SshFunction --payload '' --region us-west-2 --log-type Tail output.txt
