const dgram = require('dgram');
const server = dgram.createSocket('udp4');

const serverPort = 12345;

server.on('listening', () => {
    const address = server.address();
    console.log(`UDP Server listening on ${address.address}:${address.port}`);
});

server.on('message', (message, remote) => {
    console.log(`Received from ${remote.address}:${remote.port}: "${message.toString()}"`);
    
    // Construct a binary buffer with integer, boolean, and string values at specific byte positions
    const buffer = Buffer.alloc(16); // Allocate a buffer of 16 bytes
    // angle
    // Add an integer (4 bytes) at byte index 0
    buffer.writeInt32BE(42, 0);

    // rotate
    // Add a boolean (1 byte) at byte index 4
    buffer.writeUInt8(1, 4);

    // carrier frequency
    // Add a string at byte index 5 (use UTF-8 encoding)
    const stringValue = '70GHz';
    buffer.write(stringValue, 5, 'utf-8');

    // Send the binary buffer back to the client
    server.send(buffer, 0, buffer.length, remote.port, remote.address, (err) => {
        if (err) {
            console.error('Error sending message:', err);
        }
    });
});

server.bind(serverPort);
