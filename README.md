# AWS-Utils

Utilities for working with AWS resources.

## Resource Type Lookup

Convert an AWS ARN into a Resource Type that is supported by the [Cloud Control API](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/supported-resources.html)

> Note: Not all AWS resources are available in the Cloud Control API.

### Usage

An example of how to use the parsed ARN to pull a resource from the Cloud Control API.

> Note: This example is using the AWS v1 SDK but v2 will be similar.

```go
package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/massdriver-cloud/aws-utils/pkg/awsresourcetypes"
)

func main() {
	parsedARN, err := awsresourcetypes.Lookup("arn:aws:sqs:us-west-2:000000000000:myqueue")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsedARN)

	// This parsed ARN can then be used to pull the resource from the Cloud Control API
	cc := cloudcontrolapi.New(session.Must(session.NewSession()))
	resource, err := cc.GetResource(&cloudcontrolapi.GetResourceInput{
		Identifier: aws.String(parsedARN.ResourceId),
		TypeName:   aws.String(parsedARN.TypeName),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resource)
}
```
