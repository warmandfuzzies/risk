Framework used:
Gorilla mux

Unit tests are in main_test.go.

Run tests with the following command in the folder:
`go test`

Run the program with:
`go run main.go`

Post to the endpoint with:
`curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"state":"open","title":"I am a title","description":"I am a description"}' \
  http://localhost:8080/v1/risks`

See all risks with:
`curl http://localhost:8080/v1/risks`

Get risk with ID with:
`curl http://localhost:8080/v1/risks/{id}`


