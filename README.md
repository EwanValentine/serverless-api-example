# Serverless API Example - Golang

See write up here:

##Architecture
This repo loosely follows Uncle Bob's clean architecture. The repo itself groups functionality related to a specific domain into its own package. In this case /users.

The delivery layers include http, so you can run this repo as a standard web server. Or as a Lambda router.

The business logic is written as use cases, and we include a repository for the data layer.

##Running
As http: `$ make run-local // port 8005`.

Or deploy with serverless: `$ make deploy`.
