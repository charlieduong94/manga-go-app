package migrations

import (
  "time"
  "errors"

  "github.com/golang/glog"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

func createTable (client *dynamodb.DynamoDB, createParams *dynamodb.CreateTableInput) error {
  if createParams.TableName == nil {
    return errors.New("Table name is required")
  }

  glog.Info("Initializing table...")

  _, err := client.CreateTable(createParams)
  if err != nil {
    glog.Error(err)
    // check AND cast err if typeof awserr.Error
    if awsErr, ok := err.(awserr.Error); ok {
      // handle codes accordingly
      code := awsErr.Code()
      if (code == "ResourceInUseException") {
        glog.Info("Table already created")
      } else {
        return err
      }
    } else {
      // return the error that cannot be handled
      return err
    }
  }

  describeParams := &dynamodb.DescribeTableInput {
    TableName: createParams.TableName,
  }

  ready := false
  for ready != true {
    output, err := client.DescribeTable(describeParams)
    if err != nil {
      return err
    }

    // dereference table status for comparison
    tableStatus := *output.Table.TableStatus

    if tableStatus == "ACTIVE" {
      ready = true
    } else {
      // give the table time to complete creation
      time.Sleep(5 * time.Second)
    }
  }

  glog.Infof("Table %s is ready", *createParams.TableName)
  return nil
}

