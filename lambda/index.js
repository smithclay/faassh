'use strict';

const exec = require('child_process').exec;

exports.handler = (event, context, callback) => {
    const port = process.env.PORT || '2200';
    const child = exec(`./tiny-ssh -port ${port} -jh 172.31.24.53 -jh-user ec2-user -tunnel-port 5001`, (error) => {
        // Resolve with result of process
        //callback(error, 'Process complete!');
    });

    // Log process stdout and stderr
    child.stdout.on('data', console.log);
    child.stderr.on('data', console.error);
    // content of index.js

};
