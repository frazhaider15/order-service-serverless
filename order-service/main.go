package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Order struct {
	OrderID   int    `json:"orderID"`
	Customer  string `json:"customer"`
	ProductID int    `json:"productID"`
	Quantity  int    `json:"quantity"`
}

func CreateOrderHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request body
	var order Order
	err := json.Unmarshal([]byte(request.Body), &order)
	if err != nil {
		log.Println("Failed to parse request body:", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
	}

	// Store message in a messaging service (e.g., SQS)
	log.Println("Storing message in messaging service...")
	// ...

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Order created successfully",
	}

	return response, nil
}


func GetCustomerOrdersHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Fetch DynamoDB credentials from SSM
	ssmClient := ssm.New(session.New())
	_, err := ssmClient.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String("dynamo-db-credentials"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Println("Failed to fetch DynamoDB credentials from SSM:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// Fetch orders from DynamoDB
	dbClient := dynamodb.New(session.New())
	input := &dynamodb.ScanInput{
		TableName: aws.String("orders"),
	}
	result, err := dbClient.ScanWithContext(ctx, input)
	if err != nil {
		log.Println("Failed to fetch customer orders from DynamoDB:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	var orders []Order
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &orders)
	if err != nil {
		log.Println("Failed to unmarshal orders:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// Convert orders to JSON
	responseBody, err := json.Marshal(orders)
	if err != nil {
		log.Println("Failed to marshal response body:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return response, nil
}

func ProcessOrderHandler(ctx context.Context, order Order) error {
	// Process the order (e.g., place the order in DynamoDB, send order to messaging service)
	log.Println("Processing order:", order)

	// Place the order in DynamoDB
	_ = dynamodb.New(session.New())
	// ...

	// Send the order to a messaging service (e.g., SQS)
	// ...

	return nil
}

func UpdateStockHandler(ctx context.Context, order Order) error {
	// Update the stock in DynamoDB based on the order
	log.Println("Updating stock for order:", order)

	// Update the stock in DynamoDB
	_ = dynamodb.New(session.New())
	// ...

	return nil
}

func main() {
	lambda.Start(CreateOrderHandler)
	lambda.Start(GetCustomerOrdersHandler)
	lambda.Start(ProcessOrderHandler)
	lambda.Start(UpdateStockHandler)
}