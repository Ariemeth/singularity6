# URL Shortner

## Build
To build the application run `go build`.

## Running
After building the application to run the serer `./singularity6`

## To get a shortened url
`curl localhost:9000 -XPOST -d '{ "url": "http://www.singularity6.com" }'`

## To use the url
`curl localhost:9000/abc123`