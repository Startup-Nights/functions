# functions

Collection of serverless functions that are used mainly for the 
[startup nights website](https://github.com/Startup-Nights/website).

```sh
# initial setup of digitalocean
doctl auth init
doctl serverless install
doctl serverless connect

# set up the environment variables; these are the same as the 
# startup-nights/cli tool creates and uses to trigger the deployment via github
# actions
cp example.data.env data.env
vim data.env

# deploy the functions
bash deploy.sh
```

## Notes

Functions shouldn't use many dependencies - the resources that are used to 
build functions are very limiting. For example; to reliably get the functions 
to build, it was necessary to split the gmail / sheets functionality into two
functions.

To be able to upload files to spaces, there are specific CORS settings 
necessary. The resulting `cors.xml` is a mix of:

- https://stackoverflow.com/a/68699887
- https://stackoverflow.com/a/66558604
- https://docs.digitalocean.com/products/spaces/how-to/configure-cors/#xml
- https://github.com/digitalocean/sample-functions-golang-presigned-url/blob/main/packages/presign/url/url.go
