'use strict';

const exec = require('child_process').exec;

exports.handler = (event, context, callback) => {
    const port = process.env.PORT || '2200';
    const child = exec(`./faassh -port ${port} -jh 0.tcp.ngrok.io -jh-user csmith -jh-port 16781 -tunnel-port 5001`, (error) => {
        // Resolve with result of process
        //callback(error, 'Process complete!');
    });

    // Log process stdout and stderr
    child.stdout.on('data', console.log);
    child.stderr.on('data', console.error);
    // content of index.js

};
