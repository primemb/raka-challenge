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

	// Unmarshal to Device to access object properties
	deviceString := request.Body
	deviceStruct := Device{}
	json.Unmarshal([]byte(deviceString), &deviceStruct)

	if deviceStruct.Id == "" || deviceStruct.DeviceModel == "" || deviceStruct.Name == "" || deviceStruct.Note == "" || deviceStruct.Serial == "" {
		errorMarshalled, _ := json.Marshal(ErrorMessage{Message: "Missing required fields"})
		return events.APIGatewayProxyResponse{Body: string(errorMarshalled), StatusCode: 400}, nil
	}

	// Create new item of type item
	device := Device{
		Id:          deviceStruct.Id,
		DeviceModel: deviceStruct.DeviceModel,
		Name:        deviceStruct.Name,
		Note:        deviceStruct.Note,
		Serial:      deviceStruct.Serial,
	}

	// Marshal to dynamobb item
	av, err := dynamodbattribute.MarshalMap(device)
	if err != nil {
		fmt.Println("Error marshalling device: ", err.Error())
		errorMarshalled, _ := json.Marshal(ErrorMessage{Message: "Error store device"})
		return events.APIGatewayProxyResponse{Body: string(errorMarshalled), StatusCode: 500}, nil
	}

	tableName := os.Getenv("DYNAMODB_TABLE")

	// Build put item input
	fmt.Println("Putting Device: %v", av)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	// PutItem request
	_, err = svc.PutItem(input)

	// Checking for errors, return error
	if err != nil {
		fmt.Println("Got error calling PutItem: ", err.Error())
		errorMarshalled, _ := json.Marshal(ErrorMessage{Message: "Error store device"})
		return events.APIGatewayProxyResponse{Body: string(errorMarshalled), StatusCode: 500}, nil
	}

	// Marshal item to return
	deviceMarshalled, err := json.Marshal(device)

	fmt.Println("Returning item: ", string(deviceMarshalled))

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(deviceMarshalled), StatusCode: 201}, nil
}

func main() {
	lambda.Start(Handler)
}
