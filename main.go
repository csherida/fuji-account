package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fuji-account/internal/datastore/dynamoDB"
	"fuji-account/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Validate input paramater
	// TODO: Do a RegEx once UUID is in place
	fujiID := req.QueryStringParameters["fuji-id"]
	if _, err := strconv.Atoi(fujiID); err != nil {
		log.Printf("An invalid fujiid was given by the client.")
		return clientError(http.StatusBadRequest)
	}

	// Fetch account info from the database
	acct, err := dynamoDB.GetAccountByFujiID("1")
	if err != nil {
		return serverError(err)
	}
	if acct == nil {
		return clientError(http.StatusNotFound)
	}

	// The APIGatewayProxyResponse.Body field needs to be a string, so
	// we marshal the account record into JSON.
	js, err := json.Marshal(acct)
	if err != nil {
		return serverError(err)
	}

	// Return a response with a 200 OK status and the
	// account info wrapped in JSON as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

// This function writes a new or updated FujiAccount document
func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	acct := new(models.FujiAccount)
	err := json.Unmarshal([]byte(req.Body), acct)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// TODO: Do a RegEx once UUID is in place
	if _, err := strconv.Atoi(acct.FujiID); err != nil {
		log.Printf("An invalid fujiid was given by the client.")
		return clientError(http.StatusBadRequest)
	}
	// TODO: Validate Amazon and Apple Tokens

	err = dynamoDB.PutItem(acct)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/fujiaccount?fuji-id=%s", acct.FujiID)},
	}, nil
}

// This handles logs any error to os.Stderr and returns a
// 500 Internal Server Error response that the AWS API Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("There was a server error while handling a lambda event.")
	log.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Helper for sending responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
