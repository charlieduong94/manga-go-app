package clients

import (
  "manga-app/config"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

// creates an aws session and create a dynamo client from it
func NewDynamoDbClient () *dynamodb.DynamoDB {
  awsSession := session.Must(session.NewSession())

  conf := config.GetConfig()

  awsConfig := &aws.Config {
    Region: aws.String("us-east-1"),
  }

  // set dynamo endpoint if specified in config
  if len(conf.DynamoDbEndpoint) != 0 {
    awsConfig.Endpoint = aws.String(conf.DynamoDbEndpoint)
  }

  return dynamodb.New(awsSession, awsConfig)
}
