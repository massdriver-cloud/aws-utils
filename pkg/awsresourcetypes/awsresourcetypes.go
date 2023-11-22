package awsresourcetypes

import (
	"errors"
	"strings"

	"github.com/massdriver-cloud/aws-utils/pkg/parser"
)

type ServiceList struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	ServiceName string            `yaml:"service_name"`
	Resources   map[string]string `yaml:"resources"`
}

type ResourceType struct {
	TypeName   string
	ResourceId string
	Region     string
	FullARN    string
}

func Lookup(arn string) (ResourceType, error) {
	parsedArn, err := parser.Parse(arn)
	if err != nil {
		return ResourceType{}, err
	}

	baseResourceType := "AWS::"

	service, ok := lookupTable.Services[parsedArn.Service]
	if !ok {
		return ResourceType{}, errors.New("Service is unsuported:" + parsedArn.Service)
	}

	baseResourceType += service.ServiceName

	resource, ok := service.Resources[strings.ToLower(strings.ReplaceAll(parsedArn.Resource, "-", ""))]
	if !ok {
		return ResourceType{}, errors.New("resource type is unsupported for service:" + parsedArn.Service + ":" + parsedArn.Resource)
	}

	baseResourceType += ("::" + resource)

	return ResourceType{
		TypeName:   baseResourceType,
		ResourceId: parsedArn.ResourceId,
		Region:     parsedArn.Region,
		FullARN:    arn,
	}, nil
}
