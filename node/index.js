import net from 'node:net';
import { Readable, Transform } from 'node:stream';
import { pipeline } from 'node:stream/promises';

await main();

async function main() {
	const socket = net.createConnection({ port: 8081, host: 'localhost' }, () => {
		console.log('Connected to the server');
	});

	socket.on('data', data => {
		console.log(data.toString());
	});

	socket.on('error', err => {
		console.log(err);
	});

	const jsonTransformer = new Transform({
		objectMode: true,
		transform(chunk, encoding, callback) {
			callback(null, JSON.stringify(chunk) + '\n');
		},
	});

	const readable = Readable.from(getDataFromDb());
	readable.pipe(jsonTransformer).pipe(socket);
	// await pipeline(getDataFromDb(), jsonTransformer, socket);
	// console.log('sent all data');
}

async function* getDataFromDb() {
	let gi = 0;

	while (gi < 10) {
		const records = Array.from({ length: 1000000 }, (_, i) => ({
			id: i + gi * 100,
			name: `FUNC ${process.argv[2]} name ${i + gi * 100}`,
			price: +(Math.random() * 100).toFixed(2),
		}));

		yield records;
		gi++;
	}
}
