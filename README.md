# functions

Collection of serverless functions that are used mainly for the 
[startup nights website](https://github.com/Startup-Nights/website).

```sh
# initial setup of digitalocean
doctl auth init
doctl serverless connect
# setup the cors configuration for spaces
s3cmd setcors cors.xml s3://startupnights

# update / view tokens
go run main.go token --help

# manual deployment - data.env needs to be up to date
# automatic deployments happen via github actions on commit
# or when manually triggered
doctl serverless deploy . --env data.env 
```

## Notes

Functions shouldn't use many dependencies - the resources that are used to 
build functions are very limiting. For example; to reliably get the functions 
to build, it was necessary to split the gmail / sheets functionality into two
functions.

To be able to upload files to spaces, there are specific CORS settings 
necessary. The result in this repo is a mix of:

- https://stackoverflow.com/a/68699887
- https://stackoverflow.com/a/66558604
- https://docs.digitalocean.com/products/spaces/how-to/configure-cors/#xml
- https://github.com/digitalocean/sample-functions-golang-presigned-url/blob/main/packages/presign/url/url.go
