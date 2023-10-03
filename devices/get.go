package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"encoding/json"
	"fmt"
	"os"
)

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Creating session for client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Getting id from path parameters
	pathParamId := request.PathParameters["id"]

	fmt.Println("device id: ", pathParamId)

	// GetItem request
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(pathParamId),
			},
		},
	})

	// Checking for errors, return error
	if err != nil {
		fmt.Println(err.Error())
		errorMarshalled, _ := json.Marshal(ErrorMessage{Message: "Error getting device"})
		return events.APIGatewayProxyResponse{Body: string(errorMarshalled), StatusCode: 500}, nil
	}

	// Checking type
	if len(result.Item) == 0 {
		fmt.Println("device not found")
		errorMarshalled, _ := json.Marshal(ErrorMessage{Message: "device not found"})
		return events.APIGatewayProxyResponse{Body: string(errorMarshalled), StatusCode: 404}, nil
	}

	// Created item of type Item
	device := Device{}

	// result is of type *dynamodb.GetItemOutput
	// result.Item is of type map[string]*dynamodb.AttributeValue
	// UnmarshallMap result.item into device
	err = dynamodbattribute.UnmarshalMap(result.Item, &device)

	if err != nil {
		panic(fmt.Sprintf("Failed to UnmarshalMap result.Item: ", err))
	}

	// Marshal to type []uint8
	marshalledDevice, err := json.Marshal(device)

	// Return marshalled device
	return events.APIGatewayProxyResponse{Body: string(marshalledDevice), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
