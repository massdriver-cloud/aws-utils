package parser_test

import (
	"reflect"
	"testing"

	"github.com/massdriver-cloud/aws-resource-types/pkg/parser"
)

func TestParse(t *testing.T) {
	type test struct {
		name    string
		subject string
		want    parser.Arn
	}

	tests := []test{
		{
			name:    "Parses 7 value arns",
			subject: "arn:aws:lambda:us-west-2:000000000000:function:my-function",
			want: parser.Arn{
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
			want: parser.Arn{
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
			want: parser.Arn{
				Partition:  "aws",
				Service:    "apigateway",
				Region:     "us-west-2",
				AccountId:  "000000000000",
				Resource:   "path",
				ResourceId: "123456",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arn := parser.Parse(tc.subject)

			if !reflect.DeepEqual(tc.want, arn) {
				t.Errorf("expected result to be %+v but got %+v", tc.want, arn)
			}
		})
	}
}
