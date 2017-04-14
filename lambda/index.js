// Inspired from Apex: https://github.com/apex/apex/blob/master/shim/index.js
console.log('[shim] start function');
var child = require('child_process');

const port = process.env.PORT || '2200';
const jh = process.env.JUMP_HOST || '0.tcp.ngrok.io';
const jhPort = process.env.JUMP_HOST_PORT || '15303';
const proc = child.spawn('./faassh',
    `-port ${port} -jh ${jh} -jh-user csmith -jh-port ${jhPort} -tunnel-port 5001`.split(' '));

proc.on('error', function(err){
  console.error('[shim] error: %s', err)
  process.exit(1)
})

proc.on('exit', function(code, signal){
  console.error('[shim] exit: code=%s signal=%s', code, signal)
  process.exit(1)
});

proc.stderr.on('data', function(line){
  console.error('[faassh] data from faassh: `%s`', line)
});

proc.stdout.on('data', function(line){
  console.log('[faassh] data from faassh: `%s`', line)
});

exports.handler = (event, context, callback) => {
    context.callbackWaitsForEmptyEventLoop = false;
    // TODO: Don't kill the process, just close the session.
    setInterval(() => {
        var timeRemaining = context.getRemainingTimeInMillis();
        if (timeRemaining < 2000) {
            console.log(`Less than ${timeRemaining}ms left before timeout. Shutting down...`);
            proc.kill('SIGINT');
        }
    }, 500);
};
