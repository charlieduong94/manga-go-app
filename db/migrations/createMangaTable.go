package migrations

import (
  "manga-app/db/collections/tableNames"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

func createMangaTable (client *dynamodb.DynamoDB) error {
  tableName := aws.String(tableNames.MANGA_TABLE)

  createParams := &dynamodb.CreateTableInput{
    TableName: tableName,
    AttributeDefinitions: []*dynamodb.AttributeDefinition{
      {
        AttributeName: aws.String("id"),
        AttributeType: aws.String("S"),
      },
    },
    KeySchema: []*dynamodb.KeySchemaElement{
      {
        AttributeName: aws.String("id"),
        KeyType: aws.String("HASH"),
      },
    },
    ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
      ReadCapacityUnits: aws.Int64(10),
      WriteCapacityUnits: aws.Int64(10),
    },
  }

  return createTable(client, createParams)
}
