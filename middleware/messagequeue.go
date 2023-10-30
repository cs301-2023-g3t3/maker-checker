package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

func getAccountID() (string, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        return "", err
    }

    client := sts.NewFromConfig(cfg)

    // Use the GetCallerIdentity API to retrieve information about the calling identity
    resp, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
    if err != nil {
        return "", err
    }

    return *resp.Account, nil
}

func TriggerMessageQueueToEmail(makerEmail string, checkerEmail string) string {
   cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        return "configuration error, " + err.Error()
    }

    client := sqs.NewFromConfig(cfg)

    accountId, err := getAccountID()
    if err != nil {
        return "Unable to identify user accountId: " + err.Error()
    }

    // Provide the queue name in an environment variable
    queueName := os.Getenv("QUEUE_NAME")

    // Use GetQueueUrl to retrieve the queue URL
    gQInput := &sqs.GetQueueUrlInput{
        QueueName: aws.String(queueName),
        QueueOwnerAWSAccountId: aws.String(accountId),
    }

    result, err := GetQueueURL(context.TODO(), client, gQInput)
    if err != nil {
        return "Got an error getting the queue URL: " + err.Error()
    }

    queueURL := *result.QueueUrl

    data := map[string]interface{}{
        "makerEmail": makerEmail,
        "checkerEmail": checkerEmail,
    }

    dataJSON, err := json.Marshal(data)
    if err != nil {
        return err.Error()
    }

    sMInput := &sqs.SendMessageInput{
        DelaySeconds: 0,
        MessageBody: aws.String(string(dataJSON)),
        QueueUrl: aws.String(queueURL),
    }

    resp, err := SendMsg(context.TODO(), client, sMInput)
    if err != nil {
        return "Got an error sending the message: " + err.Error()
    }

    fmt.Println("Sent message with ID: " + *resp.MessageId)
    return ""
}
