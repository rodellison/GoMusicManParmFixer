## GoMusicManParmFixer

This is a GO rewrite of an older node/express version of MusicManParmFixer.

This is being done simply to move away from its current monolithic design, to instead 
have separate components for front end, middleware with authentication, and backend.

The front end will be a simple **React** single-page static website for presentation, 
middleware wil use API Gateway with Cognito User Pools for authentication, and the backend
will be a golang API/Lambda for handling correction inserts to the DynamoDB table.

This code in this repo provides for the serverless/lambda golang backend database handler and api gateway 
configuration. It incorporates [**Gin-gonic**](https://github.com/gin-gonic/gin), wrapped by 
[AWSLabs go api proxy](https://github.com/awslabs/aws-lambda-go-api-proxy) as a routing handler for requests. 






