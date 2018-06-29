package main

import (
	"context"

	awsLambdaGo "github.com/aws/aws-lambda-go/events"
	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	gocf "github.com/mweagle/go-cloudformation"
	"github.com/sirupsen/logrus"
)

// ASCII art: http://patorjk.com/software/taag/#p=display&v=0&f=Small&t=LAMBDA

////////////////////////////////////////////////////////////////////////////////

/*
  _      _   __  __ ___ ___   _
 | |    /_\ |  \/  | _ )   \ /_\
 | |__ / _ \| |\/| | _ \ |) / _ \
 |____/_/ \_\_|  |_|___/___/_/ \_\
*/
func sqsHandler(ctx context.Context, sqsRequest awsLambdaGo.SQSEvent) error {
	logger, _ := ctx.Value(sparta.ContextKeyLogger).(*logrus.Logger)
	logger.WithField("Event", sqsRequest).Info("SQS Event Received")
	return nil
}

/*
  __  __   _   ___ _  _
 |  \/  | /_\ |_ _| \| |
 | |\/| |/ _ \ | || .` |
 |_|  |_/_/ \_\___|_|\_|
*/
func main() {
	// 1. Create the Sparta Lambda function
	lambdaFn := sparta.HandleAWSLambda(sparta.LambdaName(sqsHandler),
		sqsHandler,
		sparta.IAMRoleDefinition{})

	// 2. Add a function decorator that provisions an SQS queue we're going to subscribe to
	sqsResourceName := "LambdaSQSFTW"
	sqsDecorator := func(serviceName string,
		lambdaResourceName string,
		lambdaResource gocf.LambdaFunction,
		resourceMetadata map[string]interface{},
		S3Bucket string,
		S3Key string,
		buildID string,
		template *gocf.Template,
		context map[string]interface{},
		logger *logrus.Logger) error {

		// Include the SQS resource in the application
		sqsResource := &gocf.SQSQueue{}
		template.AddResource(sqsResourceName, sqsResource)
		return nil
	}
	lambdaFn.Decorators = []sparta.TemplateDecoratorHandler{sparta.TemplateDecoratorHookFunc(sqsDecorator)}

	// 3. Register the lambda function as the SQS ARN EventSourceMapping subscriber.
	// See the `SQSFunctionMySqsQueue` resource in the SAM test file:
	// https://github.com/awslabs/serverless-application-model/blob/master/tests/translator/output/aws-cn/sqs.json
	lambdaFn.EventSourceMappings = append(lambdaFn.EventSourceMappings,
		&sparta.EventSourceMapping{
			EventSourceArn: gocf.GetAtt(sqsResourceName, "Arn"),
			BatchSize:      2,
		})

	// You can also "attach" to pre-existing SQS queues by using a literal or
	// gocf.String("arn:...") expression. Example:
	// lambdaFn.EventSourceMappings = append(lambdaFn.EventSourceMappings,
	// 	&sparta.EventSourceMapping{
	// 		EventSourceArn: "arn:aws:sqs:<SQS_REGION>:<SQS_ACCOUNT>:<SQS_QUEUE_NAME>",
	// 		BatchSize:      2,
	// 	})

	// Go build the stack
	sqsFunctions := []*sparta.LambdaAWSInfo{lambdaFn}
	stackName := spartaCF.UserScopedStackName("SpartaSQS")
	sparta.Main(stackName,
		"Sparta application demonstrating SQS triggering",
		sqsFunctions,
		nil,
		nil)
}
