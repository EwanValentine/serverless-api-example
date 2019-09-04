# Serverless API Example - Golang

See write up here: // Needs including

## Architecture
This repo loosely follows Uncle Bob's clean architecture. The repo itself groups functionality related to a specific domain into its own package. In this case /users.

The delivery layers include http, so you can run this repo as a standard web server. Or as a Lambda router.

The business logic is written as use cases, and we include a repository for the data layer.

## Running

### Local

If you want to run it locally with a plain http server execute the following command:

```bash
$ make run-local // port 8005
```

```
$ curl localhost:8005 
```

### Serverless

1. Apply the [example datastore](infrastructure/datastore.yml) with CloudFormation. This is required to provide an datastore with AWS DynamoDB for the Lambda Function to save data of the business logic.

2. Deploy everything else Serverless: `$ make deploy`.
