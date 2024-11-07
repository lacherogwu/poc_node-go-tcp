import net from 'net';
import JSONStream from 'JSONStream';

function handleConnection(conn, connId) {
	console.log(`Connection established: ${connId}`);

	const jsonStream = JSONStream.parse('*');

	jsonStream.on('data', record => {
		console.log('Record received', record);
	});

	jsonStream.on('error', err => {
		console.error('Error decoding JSON:', err);
		conn.end();
	});

	conn.pipe(jsonStream);

	conn.on('end', () => {
		console.log('Connection closed', connId);
	});

	conn.on('error', err => {
		console.error('Connection error:', err);
	});
}

const server = net.createServer(conn => {
	const connId = Math.floor(Math.random() * 1000);
	handleConnection(conn, connId);
});

server.listen(8081, () => {
	console.log('Server listening on port 8081');
});
