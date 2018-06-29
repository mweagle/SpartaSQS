## SpartaSQS

Sample [Sparta](https://gosparta.io) application that demonstrates:

  1. Defining an AWS Lambda function that accepts a [SQSEvent](https://godoc.org/github.com/aws/aws-lambda-go/events#SQSEvent) request
  1. Provisioning an SQS Queue using a [TemplateDecorator](https://godoc.org/github.com/mweagle/Sparta#TemplateDecoratorHandler)
  1. Subscribing the Lambda function to SQS messages via an [EventSourceMapping](https://godoc.org/github.com/mweagle/Sparta#EventSourceMapping)

