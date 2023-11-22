package awsresourcetypes_test

import (
	"testing"

	"github.com/massdriver-cloud/aws-utils/pkg/awsresourcetypes"
)

func TestLookup(t *testing.T) {
	type test struct {
		name    string
		subject string
		want    awsresourcetypes.ResourceType
	}

	tests := []test{
		{
			name:    "Handles 7 value arns",
			subject: "arn:aws:lambda:us-west-2:000000000000:function:my-function",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::Lambda::Function",
				ResourceId: "my-function",
				Region:     "us-west-2",
				FullARN:    "arn:aws:lambda:us-west-2:000000000000:function:my-function",
			},
		},
		{
			name:    "Handles path style resource ids",
			subject: "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::ApiGateway::RestApi",
				ResourceId: "wymjfx3iie",
				Region:     "us-west-2",
				FullARN:    "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie",
			},
		},
		{
			name:    "Handles relative ids",
			subject: "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie/resource/123456",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::ApiGateway::Resource",
				ResourceId: "123456",
				Region:     "us-west-2",
				FullARN:    "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie/resource/123456",
			},
		},
		{
			name:    "Handles S3",
			subject: "arn:aws:s3::000000000000:myBucket",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::S3::Bucket",
				ResourceId: "myBucket",
				Region:     "",
				FullARN:    "arn:aws:s3::000000000000:myBucket",
			},
		},
		{
			name:    "Handles SQS",
			subject: "arn:aws:sqs:us-west-2:000000000000:myqueue",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::SQS::Queue",
				ResourceId: "myqueue",
				Region:     "us-west-2",
				FullARN:    "arn:aws:sqs:us-west-2:000000000000:myqueue",
			},
		},
		{
			name:    "Handles SnS",
			subject: "arn:aws:sns:us-west-2:000000000000:mytopic",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::SNS::Topic",
				ResourceId: "mytopic",
				Region:     "us-west-2",
				FullARN:    "arn:aws:sns:us-west-2:000000000000:mytopic",
			},
		},
		{
			name:    "Handles Resources at the beginning",
			subject: "arn:aws:ecs:us-west-2:083014189801:task-definition/mclecs-tes-testces-0plc:1",
			want: awsresourcetypes.ResourceType{
				TypeName:   "AWS::ECS::TaskDefinition",
				ResourceId: "mclecs-tes-testces-0plc",
				Region:     "us-west-2",
				FullARN:    "arn:aws:ecs:us-west-2:083014189801:task-definition/mclecs-tes-testces-0plc:1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resourceType, err := awsresourcetypes.Lookup(tc.subject)

			if err != nil {
				t.Fatal(err)
			}

			if tc.want != resourceType {
				t.Errorf("expected %s but got %s", tc.want, resourceType)
			}
		})
	}
}
