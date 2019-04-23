# faassh
### simple go SSH server designed for running in cloud functions

![image](https://cloud.githubusercontent.com/assets/27153/25602411/819d0b02-2ea8-11e7-9f64-157226b2d4cb.png)

This is just for fun. It's a simple SSH server and tunnel-er that allows you to SSH into a running lambda function—until it times out.

Developed for my [dotScale](https://dotscale.io) 2017 talk, "Searching for the Server in Serverless". Slides [here](http://speakerdeck.com/smithclay/searching-for-the-server-in-serverless).

## building

This project uses the [Serverless Application Model](https://aws.amazon.com/serverless/sam/) for packaging and deploying.

```sh
   $ sam build
   $ sam package --s3-bucket <yourbucket> > packaged.yaml
   $ sam deploy --template-file packaged.yaml --stack-name <yourstack> --capabilities CAPABILITY_IAM
```

## usage

```
   faassh -i ./path_to_private_rsa_host_key -p port_number
```

## example

See the example node.js lambda function in the `lambda/` directory.

* Generate RSA keys for the Lambda function and bundle inside the `lambda` directory (`ssh-keygen -t rsa -f ./id_rsa`)
* Set the envionment variables to point to your SSH jump host with the correct username.

If you'd like to test it on your local laptop that's behind (hopefully) a NAT/firewall, I like the TCP forwarding available on [ngrok](https://ngrok.com/). You can create a tunnel to your local SSH server for the other end of the tunnel endpoint, you just run: `ngrok tcp 22`.

## other interesting/related projects

* [lambdash](https://github.com/alestic/lambdash) - another approach for running commands in Lambda
* [awslambdaproxy](https://github.com/dan-v/awslambdaproxy) - An AWS Lambda powered HTTP/SOCKS web proxy

## todo

- better authentication support
- other cloud providers
- connection cleanup
- terraform/cloudformation helper
- multiple connections
- tests and docs :)
