package awsresourcetypes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/massdriver-cloud/aws-resource-types/pkg/parser"
	"gopkg.in/yaml.v3"
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
}

func Lookup(arn string) (ResourceType, error) {
	table, err := createLookupTable()

	if err != nil {
		panic(err)
	}

	parsedArn := parser.Parse(arn)
	fmt.Println(parsedArn)
	baseResourceType := "AWS::"

	service, ok := table.Services[parsedArn.Service]

	if !ok {
		return ResourceType{}, errors.New("Service is unsuported")
	}

	baseResourceType += service.ServiceName

	resource, ok := service.Resources[parsedArn.Resource]

	if !ok {
		return ResourceType{}, errors.New("Resource Type is unsupported for service")
	}

	baseResourceType += ("::" + resource)

	return ResourceType{
		TypeName:   baseResourceType,
		ResourceId: parsedArn.ResourceId,
	}, nil
}

// I hope we can move this to compile time.
func createLookupTable() (ServiceList, error) {
	var serviceList ServiceList

	filename, _ := filepath.Abs("./resource_types_lookup.yaml")
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		return serviceList, err
	}

	err = yaml.Unmarshal(yamlFile, &serviceList)

	if err != nil {
		return serviceList, err
	}

	return serviceList, nil
}
