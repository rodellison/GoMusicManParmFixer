package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	muxadapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/rodellison/gomusicmanparmfixer/common"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	muxLambda          *muxadapter.GorillaMuxAdapter
	ProcessDataHandler func(w http.ResponseWriter, req *http.Request)
)

type ResponseOutput struct {
	Message string   `json:"message"`
}

func init() {

	ProcessDataHandler = ProcessDatabase

}

func extractRequestData(req *http.Request) (common.RequestParmData, error) {

	//	myHeader := c.GetHeader("Authorization")

	//Retrieve the incoming Request Body, and unmarshal it into a variable we can use
	x, _ := ioutil.ReadAll(req.Body)
	log.Println("Request Body from Context:")
	log.Println(string(x))

	var reqParmData common.RequestParmData
	json.Unmarshal(x, &reqParmData)

	return reqParmData, nil

}

func ProcessDatabase (w http.ResponseWriter, req *http.Request) {

	status := 200
	message := "Processed request successfully!"

	myRequestBodyParmData, err := extractRequestData(req)

	if err == nil {
		//Need to switch the key to lower case
		myRequestBodyParmData.SongKickInvalidParmData = strings.ToLower(myRequestBodyParmData.SongKickInvalidParmData)
		//Process the request data here
		dbErr := common.PutDBParmData(req.Context(), myRequestBodyParmData)
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

	myResponse := ResponseOutput{
		Message: message,
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(myResponse)

}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if muxLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("GorillaMux cold start")
		r := mux.NewRouter()

		r.HandleFunc("/parmdata", ProcessDataHandler).Methods("POST")

		muxLambda = muxadapter.New(r)
	}

	return muxLambda.ProxyWithContext(ctx, req)

}

func main() {
	lambda.Start(Handler)
}
