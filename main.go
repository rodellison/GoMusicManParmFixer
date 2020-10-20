package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/rodellison/gomusicmanparmfixer/common"
	"io/ioutil"
	"log"
	"strings"
)

var (
	ginLambda          *ginadapter.GinLambda
	ProcessDataHandler func(*gin.Context)
)

func init() {

	ProcessDataHandler = ProcessDatabase

}

func extractRequestData(c *gin.Context) (common.RequestParmData, error) {

	//	myHeader := c.GetHeader("Authorization")

	//Retrieve the incoming Request Body, and unmarshal it into a variable we can use
	x, _ := ioutil.ReadAll(c.Request.Body)
	log.Println("Request Body from Context:")
	log.Println(string(x))

	var reqParmData common.RequestParmData
	json.Unmarshal(x, &reqParmData)

	return reqParmData, nil

}

func ProcessDatabase(c *gin.Context) {

	status := 200
	message := "Processed request successfully!"
	myRequestBodyParmData, err := extractRequestData(c)

	if err == nil {
		//Need to switch the key to lower case
		myRequestBodyParmData.SongKickInvalidParmData = strings.ToLower(myRequestBodyParmData.SongKickInvalidParmData)
		//Process the request data here
		dbErr := common.PutDBParmData(c.Request.Context(), myRequestBodyParmData)
		if dbErr != nil {
			message = dbErr.Error()
			status = 400
		} else {
			status = 200
			message = "Processed request successfully!"
		}
	} else {
		status = 400
		message = err.Error()
	}

	c.JSON(status, gin.H{
		"message": message,
	})

}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()

		//Trying out the ability to distinguish API versions 
		dev := r.Group("/V1")
		{
			dev.POST("/parmdata", ProcessDataHandler)
		}

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
//	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
