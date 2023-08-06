# functions

Collection of serverless functions that are used maily for the 
[startup nights website](https://github.com/Startup-Nights/website).

- `./package/website/gmail`: send an email to someone
- `./package/website/sheets`: write data to a sheet or read from it
- `./package/website/spaces`: create a presigned url to upload images to 
  digitalocean spaces

## Setup and deployment
 
```sh
# initial setup of digitalocean
doctl auth init
doctl serverless connect
# setup the cors configuration for spaces
s3cmd setcors cors.xml s3://startupnights

# create credentials for gmail / sheets API and update the env file
go run cmd/token/main.go

# deploy and update the the functions as defined in 'project.yml'
doctl serverless deploy . --env data.env 
```

## Notes

Deploying / building the functions is not very consistent - there is often a 
failure. Upon trying again, it often just vanishes. No solution yet; increasing
both memory and timeout didn't solve it yet. The only option left is to try and
reduce the imports - and split functions into multiple functions

To be able to upload files to spaces, there are specific CORS settings 
necessary. The result in this repo is a mix of:

- https://stackoverflow.com/a/68699887
- https://stackoverflow.com/a/66558604
- https://docs.digitalocean.com/products/spaces/how-to/configure-cors/#xml
- https://github.com/digitalocean/sample-functions-golang-presigned-url/blob/main/packages/presign/url/url.go
