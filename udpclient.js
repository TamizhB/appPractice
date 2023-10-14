const dgram = require('dgram');
const client = dgram.createSocket('udp4');

const serverAddress = 'localhost';
const serverPort = 12345;

const message = 'getconfig';

const messageBuffer = Buffer.from(message);

client.send(messageBuffer, 0, messageBuffer.length, serverPort, serverAddress, (err) => {
    if (err) {
        console.error('Error sending message:', err);
        client.close();
    } else {
        console.log(`Sent to ${serverAddress}:${serverPort}: "${message}"`);
    }
});

client.on('message', (response, remote) => {
    console.log(`Received from ${remote.address}:${remote.port}: "${response.toString()}"`);

    // Parse the binary buffer received from the server
    const angle = response.readInt32BE(0); // Read the integer at byte index 0
    const rotate = response.readUInt8(4) === 1; // Read the boolean at byte index 4
    const frequency = response.toString('utf-8', 5); // Read the string starting at byte index 5

    console.log('Parsed Data:');
    console.log(`angle: ${angle}`);
    console.log(`rotate: ${rotate}`);
    console.log(`frequency: ${frequency}`);

    client.close();
});

client.on('error', (err) => {
    console.error('UDP client error:', err);
});

client.on('close', () => {
    console.log('UDP client closed');
});
