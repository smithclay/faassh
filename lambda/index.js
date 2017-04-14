'use strict';

console.log('start function');

const exec = require('child_process').exec;

exports.handler = (event, context, callback) => {
    const port = process.env.PORT || '2200';
    const jh = process.env.JUMP_HOST || '0.tcp.ngrok.io';
    const jhPort = process.env.JUMP_HOST_PORT || '15303';
    const child = exec(`./faassh -port ${port} -jh ${jh} -jh-user csmith -jh-port ${jhPort} -tunnel-port 5001`);

    setInterval(() => {
        var timeRemaining = context.getRemainingTimeInMillis();
        if (timeRemaining < 2000) {
            console.log(`Less than ${timeRemaining}ms left before timeout. Shutting down...`);
            child.kill('SIGINT');
        }
    }, 500);

    child.on('error', (error) => {
        console.log(`error recieved from child process: ${error}`);
        return callback(error, null);
    });

    // TODO: this doesn't seem to get called
    child.on('exit', (code, signal) => {
        console.log(`exit signal recieved from child process: ${code} ${signal}`);
        if (code !== 0) {
            return callback(code, null);
        }
        callback(null, code);
    });

    // Log process stdout and stderr
    child.stdout.on('data', console.log);
    child.stderr.on('data', console.error);
};
