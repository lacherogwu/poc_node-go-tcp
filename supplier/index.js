import http from 'node:http';
import fs from 'node:fs';

const htmlFile = fs.readFileSync('./index.html', 'utf8');
let totalReq = 0;
const server = http.createServer(async (req, res) => {
	let reqNum = ++totalReq;
	console.log(`Request received #${reqNum}`);
	await new Promise(resolve => setTimeout(resolve, Math.random() * 5000));

	res.writeHead(200, { 'Content-Type': 'text/plain' });
	res.end(htmlFile);
});

server.listen(3000, () => {
	console.log('Server is listening on port 3000');
});
