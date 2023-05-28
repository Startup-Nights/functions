# functions

Collection of serverless functions that are used maily for the 
[startup nights website](https://github.com/Startup-Nights/website). In the 
future they might get moved to the website itself since vercel has functions as
well.

```sh
# initial setup of digitalocean
doctl auth init
doctl serverless connect

# create credentials for gmail / sheets API and update the env file
go run cmd/token/main.go

# deploy and update the the functions as defined in 'project.yml'
doctl serverless deploy . --env data.env 
```

## Notes

To be able to upload files to spaces, there are specific CORS settings 
necessary. The result in this repo is a mix of:

- https://stackoverflow.com/a/68699887
- https://stackoverflow.com/a/66558604
- https://docs.digitalocean.com/products/spaces/how-to/configure-cors/#xml
- https://github.com/digitalocean/sample-functions-golang-presigned-url/blob/main/packages/presign/url/url.go