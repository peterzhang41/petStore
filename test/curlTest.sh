#AddPet
curl -X POST "http://127.0.0.1:8080/v2/pet" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"id\": 0, \"category\": { \"id\": 0, \"name\": \"string\" }, \"name\": \"doggie\", \"photoUrls\": [ \"string\" ], \"tags\": [ { \"id\": 0, \"name\": \"string\" } ], \"status\": \"available\"}"

#UpdatePet
curl -X PUT "http://127.0.0.1:8080/v2/pet" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"id\": 0, \"category\": { \"id\": 0, \"name\": \"string\" }, \"name\": \"doggie\", \"photoUrls\": [ \"string\" ], \"tags\": [ { \"id\": 0, \"name\": \"string\" } ], \"status\": \"available\"}"

#GetPetById
curl -X GET "http://127.0.0.1:8080/v2/pet/1" -H "accept: application/json"

#UpdatePetWithForm
curl -X POST "http://127.0.0.1:8080/v2/pet/1" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "name=4&status=6"

#DeletePet
curl -X DELETE "http://127.0.0.1:8080/v2/pet/2" -H "accept: application/json"

#UploadImage
curl -X POST "http://127.0.0.1:8080/v2/pet/3/uploadImage" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "additionalMetadata=pdf" -F "file=@upload_test_file.log;type=text/x-log"

#FindByStatus
curl -X GET "http://127.0.0.1:8080/v2/pet/findByStatus?status=available&status=pending&status=sold" -H "accept: application/json"

#FindByTags
curl -X GET "http://127.0.0.1:8080/v2/pet/findByTags?tags=test1&tags=test1&tags=test2&tags=test3&tags=test4" -H "accept: application/json"