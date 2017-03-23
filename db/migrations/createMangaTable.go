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
    AttributeDefinitions: []*dynamodb.AttributeDefinition {
      // type attr is used as a hash key for indexes
      {
        AttributeName: aws.String("lang"),
        AttributeType: aws.String("S"),
      },
      {
        AttributeName: aws.String("id"),
        AttributeType: aws.String("S"),
      },
      {
        AttributeName: aws.String("lastChapterDate"),
        AttributeType: aws.String("N"),
      },
      {
        AttributeName: aws.String("hits"),
        AttributeType: aws.String("N"),
      },
    },
    KeySchema: []*dynamodb.KeySchemaElement {
      {
        AttributeName: aws.String("lang"),
        KeyType: aws.String("HASH"),
      },
      {
        AttributeName: aws.String("id"),
        KeyType: aws.String("RANGE"),
      },
    },
    ProvisionedThroughput: &dynamodb.ProvisionedThroughput {
      ReadCapacityUnits: aws.Int64(10),
      WriteCapacityUnits: aws.Int64(10),
    },
    LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex {
      {
        IndexName: aws.String("lastChapterIndex"),
        KeySchema: []*dynamodb.KeySchemaElement {
          {
            AttributeName: aws.String("lang"),
            KeyType: aws.String("HASH"),
          },
          {
            AttributeName: aws.String("lastChapterDate"),
            KeyType: aws.String("RANGE"),
          },
        },
        Projection: &dynamodb.Projection {
          ProjectionType: aws.String("ALL"),
        },
      },
      {
        IndexName: aws.String("hitsIndex"),
        KeySchema: []*dynamodb.KeySchemaElement {
          {
            AttributeName: aws.String("lang"),
            KeyType: aws.String("HASH"),
          },
          {
            AttributeName: aws.String("hits"),
            KeyType: aws.String("RANGE"),
          },
        },
        Projection: &dynamodb.Projection {
          ProjectionType: aws.String("ALL"),
        },
      },
    },
    StreamSpecification: &dynamodb.StreamSpecification {
      StreamEnabled: aws.Bool(true),
      StreamViewType: aws.String("NEW_IMAGE"),
    },
  }

  return createTable(client, createParams)
}
