package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rodellison/gomusicmanparmfixer/common"
	"github.com/rodellison/gomusicmanparmfixer/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	//IMPORTANT!! - for the test to use our mocked response below, we have to make sure to set the client to
	//be the mocked client, which will use the overridden versions of the function that makes calls
	//During testing, we'll override the endpoint to ensure testing against local DynamoDB Docker image
//	cfg := aws.Config{
//		Endpoint:   aws.String("http://localhost:8000"),
//		Region:     aws.String("us-east-1"),
//		MaxRetries: aws.Int(3),
//	}

	//Get Session, credentials
//	sess := session.Must(session.NewSessionWithOptions(session.Options{
//		SharedConfigState: session.SharedConfigEnable,
//	}))

//	common.DynamoDBSvcClient = dynamodb.New(sess, &cfg) //use this one for actual dB interaction - test or prod
	common.DynamoDBIfaceClient = &mocks.MockDynamoDBSvcClient{}
	//common.SNSIfaceClient = &mocks.MockSNSSvcClient{}


}

func TestHandlerCanInsertDynamoDBRequest(t *testing.T) {


	expectedResult := `{StatusCode: 200,Body:"Thanks"}`

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  interface{}
		err     error
	}{
		{
			request: nil,
			expect:  expectedResult,
			err:     nil,
		},
	}



	// build mock DynamoDB put
	//mocks.MockDynamoPutItem = func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	//	//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
	//	fmt.Println("Mock DynamoDB Put Item called")
	//	return &dynamodb.PutItemOutput{
	//		Attributes:            nil,
	//		ConsumedCapacity:      nil,
	//		ItemCollectionMetrics: nil,
	//	}, nil
	//}

	//mocks.MockDoSNSPublish = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS Publish called with info: " + *input.Message)
	//	return &sns.PublishOutput{}, nil
	//}
	//mocks.MockDoSNSPublishWithContext = func(ctx aws.Context, input *sns.PublishInput, options ...request.Option) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS PublishWithContext called with info: " + *input.Message)
	//	return &sns.PublishOutput{}, nil
	//}

	var ctx context.Context
	ctx = context.Background()


	for _, test := range tests {
		response, err := Handler(ctx, test.request)
		assert.IsType(t, test.err, err)
		assert.Contains(t, response.StatusCode, 200 )
	}

}
