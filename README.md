# PetStore Restful API demo

A REST API microservice written in GO

## Prerequisite

Go or Docker installed

## How to run

- Running from go
 
    ```
        go get github.com/peterzhang41/petStore  
        go run ~/go/src/github.com/peterzhang41/petStore
    ```

- Running from docker
    ```
        docker build -t go-docker-api-service .
        docker run -d -p 8080:8080 --name pet-store-api go-docker-api-service
        docker logs pet-store-api -f
    ```

## Notice
-  All pet endpoints are fully operational, other operations are not implemented.
-  Due to time limit, data is stored in memory, however a SQL installation file in the DB folder for future use.
-  A Dockerfile is written and added, which makes this service self-contained and highly scalable.
-  Models are auto-generated from swagger auto-gen tool.

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




