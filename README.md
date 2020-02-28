# PetStore Restful API demo

A REST API microservice written in GO

## Prerequisite

- GO installed [https://golang.org/]
- localhost:8080 is not in use

## How to run

```go get github.com/peterzhang41/petStore```

```go run ~/go/src/github.com/peterzhang41/petStore```

## Notice
-  All pet endpoints are fully operational, other operations are not coded
-  Due to time limit, data is stored in memory (a hash map), however a SQL installation file in the DB folder for future use
-  Uploaded files will be in project root folder
-  Some error handling code is a little redundant which could be refactored
-  Curl test cases are not real automation tests. Unit tests could be added
-  A Dockerfile could be written and added to make this service self-contained
-  Models are auto-generated from swagger auto-gen tool

## TESTING

```
cd test
./curlTest.sh

```

## API Requirement

https://petstore.swagger.io




