# POC

node process pull from db 20k records
then sending all the records to go process via tcp
go process will insert all the records to channel and will have x workers to process the records
each record that have been processed will be sent back to node process via tcp
node will continue to process it and insert the final result to db

## Plan

Node pulls all info that needs to be scraped from db
then organize it in a uniform way on how to scrape, using lambda, using our tcp server
or its always happens through our tcp server

Options:

- lambda
- proxy:residential
- proxy:datacenter
- proxy:mobile

// data

```json
{
	"url": "https://www.google.com",
	"method": "GET",
	// "headers": {
	// 	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	// },
	"headers": "http://generate-champs-headers",
	"body": "anything",
	"retries": 3,
	"timeout": 10000,
	"retryConditions": [],
	// "proxy": "proxy:residential-dicks-countries"
	"proxy": "lambda"
}
```
