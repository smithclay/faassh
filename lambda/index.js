'use strict';

const exec = require('child_process').exec;

exports.handler = (event, context, callback) => {
    const port = process.env.PORT || '2200';
    const child = exec(`./faassh -port ${port} -jh 0.tcp.ngrok.io -jh-user csmith -jh-port 16781 -tunnel-port 5001`);

    setInterval(() => {
        var timeRemaining = context.getRemainingTimeInMillis();
        if (timeRemaining < 2000) {
            console.log(`Less than ${timeRemaining}ms left before timeout. Shutting down...`);
            child.kill('SIGINT');
        }
    }, 1000);

    child.on('error', (error) => {
        return callback(error, null);
    });

    // TODO: this doesn't seem to get called
    child.on('exit', (code, signal) => {
        if (code !== 0) {
            return callback(code, null);
        }
        callback(null, code);
    });

    // Log process stdout and stderr
    child.stdout.on('data', console.log);
    child.stderr.on('data', console.error);
};
