package parser_test

import (
	"testing"

	"github.com/massdriver-cloud/aws-utils/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type test struct {
		name    string
		subject string
		want    *parser.Arn
		err     bool
	}

	tests := []test{
		{
			name:    "Parses 7 value arns",
			subject: "arn:aws:lambda:us-west-2:000000000000:function:my-function",
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "lambda",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "function",
				ResourceId: "my-function",
			},
		},
		{
			name:    "Parses path style resource ids",
			subject: "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie",
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "apigateway",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "restapi",
				ResourceId: "wymjfx3iie",
			},
		},
		{
			name:    "Parses Relative Ids",
			subject: "arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie/path/123456",
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "apigateway",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "path",
				ResourceId: "123456",
			},
		},
		{
			name:    "Parses generic resource-type/ID ARN",
			subject: "arn:aws:ec2:us-east-1:000000000000:subnet/subnet-wymjfx3iie",
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "ec2",
				Region:     "us-east-1",
				AccountId:  "000000000000",
				Resource:   "subnet",
				ResourceId: "subnet-wymjfx3iie",
			},
		},
		{
			name:    "Parses ECS style",
			subject: "arn:aws:ecs:us-west-2:000000000000:task-definition/wymjfx3iie:1", // The revision number is dropped on parse but the api pulls latest if not included
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "ecs",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "task-definition",
				ResourceId: "wymjfx3iie",
			},
		},
		{
			name:    "Parses EKS style",
			subject: "arn:aws:eks:us-west-2:000000000000:nodegroup/cool-cluster-name1234/cool-cluster-name1234-shared-2/123456-f14e-fb7e-91fc-8dc78e043e15",
			want: &parser.Arn{
				Partition:  "aws",
				Service:    "eks",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "nodegroup",
				ResourceId: "cool-cluster-name1234/cool-cluster-name1234-shared-2/123456-f14e-fb7e-91fc-8dc78e043e15",
			},
		},
		{
			name:    "Fails on invalid prefix",
			subject: "not-an-arn:aws:apigateway:us-west-2:000000000000:/restapis/wymjfx3iie/path/123456",
			want:    nil,
			err:     true,
		},
		{
			name:    "Fails on insufficient parts on split",
			subject: "arn:aws:apigateway:us-west-2",
			want:    nil,
			err:     true,
		},
		{
			name:    "Fails on insufficient parts on resourceID split",
			subject: "arn:aws:ec2:us-east-1:000000000000:subnet-06676afb0ac48da5d",
			want:    nil,
			err:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			arn, err := parser.Parse(tc.subject)

			assert.Equal(tc.want, arn)

			if tc.err {
				assert.NotNil(err, "expected an error")
			} else {
				// This should get caught above in the Equal check
				assert.Nil(err, "expected no error")
			}
		})
	}
}
