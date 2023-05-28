# functions

## Functions

```sh
# initial setup
doctl auth init
doctl serverless connect

# deploy and update the the functions
doctl serverless deploy . --env data.env 
```

## CORS

To be able to upload files to spaces, there are specific CORS settings 
necessary. The result in this repo is a mix of:

- https://stackoverflow.com/a/68699887
- https://stackoverflow.com/a/66558604
- https://docs.digitalocean.com/products/spaces/how-to/configure-cors/#xml