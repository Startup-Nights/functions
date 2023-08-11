# functions

Collection of serverless functions that are used maily for the 
[startup nights website](https://github.com/Startup-Nights/website).

- `./package/website/gmail`: send an email to someone
- `./package/website/sheets`: write data to a sheet or read from it
- `./package/website/spaces`: create a presigned url to upload images to 
  digitalocean spaces

The interfaces / requests:

```go 
// gmail
type Request struct {
	Recipient string `json:"recipient"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

// sheets
type Request struct {
	ID    string   `json:"id"`
	Range string   `json:"range"`
	Data  []string `json:"data"`
}

// spaces
type Request struct {
	Filename string `json:"filename"`
}
```

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
