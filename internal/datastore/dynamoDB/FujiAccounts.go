package dynamoDB

import (
	"context"
	"fmt"
	"fuji-account/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Fetch a Fuji account from the database using the FujiID
func GetAccountByFujiID(fujiID string) (*models.FujiAccount, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	result, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("FujiAccounts"),
		Key: map[string]types.AttributeValue{
			"FujiID": &types.AttributeValueMemberS{Value: fujiID},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(result.Item)

	// Map the result to FujiAccount
	acct := new(models.FujiAccount)
	err = attributevalue.UnmarshalMap(result.Item, acct)
	if err != nil {
		return nil, err
	}

	return acct, nil
}

// Fetch a Fuji account from the database using the Amazon Token
func GetAccountByAmazonToken(amazonToken string) (*models.FujiAccount, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	result, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("FujiAccounts"),
		IndexName:              aws.String("AmazonToken-index"),
		KeyConditionExpression: aws.String("AmazonToken = :AmazonToken"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":AmazonToken": &types.AttributeValueMemberS{Value: amazonToken},
			//":gsi1sk": &types.AttributeValueMemberN{Value: "20150101"},
		},
	})
	if err != nil {
		panic(err)
	}

	// Map the result to FujiAccount
	acct := new(models.FujiAccount)
	//TODO: Handle if more than 1
	err = attributevalue.UnmarshalMap(result.Items[0], acct)
	if err != nil {
		return nil, err
	}

	return acct, nil
}

// Add a new Fuji account to DynamoDB.
func PutItem(acct *models.FujiAccount) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("FujiAccount"),
		Item: map[string]types.AttributeValue{
			"FujiID":      &types.AttributeValueMemberS{Value: acct.FujiID},
			"AmazonToken": &types.AttributeValueMemberS{Value: acct.AmazonToken},
			"AppleToken":  &types.AttributeValueMemberS{Value: acct.AppleToken},
		},
	})

	fmt.Println(out.Attributes)
	return err
}
