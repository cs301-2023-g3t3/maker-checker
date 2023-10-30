package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSSendMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSSendMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func SendMsg(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}

func parseAccountIDFromARN(arn string) (string, error) {
	// ARN format: arn:aws:iam::123456789012:user/username
	// Extract account ID (digits between the last two colons)
	arnParts := strings.Split(arn, ":")
	if len(arnParts) < 5 {
		return "", fmt.Errorf("Invalid ARN format")
	}

	return arnParts[4], nil
}

func getAccountID() (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		return "", err
	}

	client := iam.NewFromConfig(cfg)

	// Use the GetUser API to retrieve information about the current user
	// This will include the user's ARN, which includes the account ID
	userInfo, err := client.GetUser(context.TODO(), &iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	// Parse the user's ARN to extract the account ID
	accountID, err := parseAccountIDFromARN(*userInfo.User.Arn)
	if err != nil {
		return "", err
	}

	return accountID, nil
}

func TriggerMessageQueueToEmail(makerEmail string, checkerEmail string) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        panic("configuration error, " + err.Error())
    }

    client := sqs.NewFromConfig(cfg)

    // Get URL of queue (if you need this, provide the queue name in an environment variable)

    accountId, err := getAccountID()
    if err != nil {
        panic(err.Error())
    }

    queueURL := fmt.Sprintf("https://sqs.ap-southeast-1.amazonaws.com/%v/%v", accountId, os.Getenv("QUEUE_NAME"))

    data := map[string]interface{}{
        "makerEmail": makerEmail,
        "checkerEmail": checkerEmail,
    }

    dataJSON, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }

    sMInput := &sqs.SendMessageInput{
        DelaySeconds: 0,
        MessageGroupId: aws.String("makerchecker"),
        MessageBody: aws.String(string(dataJSON)),
        QueueUrl:    aws.String(queueURL),
    }

    resp, err := SendMsg(context.TODO(), client, sMInput)
    if err != nil {
        fmt.Println("Got an error sending the message:")
        fmt.Println(err)
        return
    }

    fmt.Println("Sent message with ID: " + *resp.MessageId)
}
