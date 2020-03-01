# PetStore Restful API demo

A REST API microservice written in GO

## Prerequisite

- GO installed [https://golang.org/]
- localhost:8080 is not in use

## How to run

```go get github.com/peterzhang41/petStore```

```go run ~/go/src/github.com/peterzhang41/petStore```

## Notice
-  All pet endpoints are fully operational, other operations are not implemented
-  Due to time limit, data is stored in memory, however a SQL installation file in the DB folder for future use
-  A Dockerfile could be written and added to make this service self-contained
-  Models are auto-generated from swagger auto-gen tool

## TESTING
Pet endpoint integration tests are written.  
Results are all passed and 85% statements coverage

```
go test ./... -cover -v
```
curlTest.sh in the test folder is a tool for quick access checking
```
test/curlTest.sh
```

## API Requirement

https://petstore.swagger.io




