# URL Shortner

## Build
To build the application run `go build`.

## Running
After building the application to run the serer `./singularity6`

## To get a shortened url
`curl localhost:9000 -XPOST -d '{ "url": "http://www.singularity6.com" }'`

## To use the url
`curl localhost:9000/abc123`

## Health
To check if the service is running the health endpoint can be checked
`curl localhost:9000/healthz -X POST`