package dbClients

import (
  "manga-app/config"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

// creates an aws session and create a dynamo client from it
func NewDynamoDBClient () *dynamodb.DynamoDB {
  awsSession := session.Must(session.NewSession())

  conf := config.GetConfig()

  awsConfig := &aws.Config {
    Region: aws.String("us-east-1"),
  }

  if len(conf.DynamoDBEndpoint) != 0 {
    awsConfig.Endpoint = aws.String("http://localhost:8000")
  }

  return dynamodb.New(awsSession, awsConfig)
}
