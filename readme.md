# faassh
### simple go SSH server designed for running in cloud functions

![image](https://cloud.githubusercontent.com/assets/27153/25602411/819d0b02-2ea8-11e7-9f64-157226b2d4cb.png)

This is just for fun. It's a simple SSH server and tunnel-er that allows you to SSH into a running lambda function—until it times out.

## usage

```
   faassh -i ./path_to_private_rsa_host_key -p port_number
```

## example

See the example node.js lambda function in the `lambda/` directory.

* Generate RSA keys for the Lambda function and bundle inside the `lambda` directory (`ssh-keygen -t rsa -f ./id_rsa`)
* Set the envionment variables to point to your SSH jump host with the correct username.

If you'd like to test it on your local laptop that's behind (hopefully) a NAT/firewall, I like the TCP forwarding available on [ngrok](https://ngrok.com/).

## todo

- better authentication support
- other cloud providers
- connection cleanup
- terraform/cloudformation helper
- multiple connections
- tests and docs :)
