package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Arn struct {
	Partition  string
	Service    string
	Region     string
	AccountId  string
	Resource   string
	ResourceId string
}

var partitionIndex = 1
var serviceIndex = 2
var regionIndex = 3
var accountIndex = 4
var resourceIndex = 5
var resourceIdIndex = 6
var standardLength = 6

var ServiceBaseTypes = map[string]string{
	"s3":  "bucket",
	"sns": "topic",
	"sqs": "queue",
}

// resourceBases are resources that have their type at the beginning
// of the resource and any number of / and : after that with the section after
// the first / being the ID
var resourceBases = map[string]struct{}{
	"ecs":                    {},
	"eks":                    {},
	"ecr":                    {},
	"elasticloadbalancingv2": {},
	"events":                 {},
	"sagemaker":              {},
}

var serviceMapping = map[string]string{
	"acm":                  "acmpca",
	"catalog":              "servicecatalog",
	"elasticfilesystem":    "efs",
	"elasticloadbalancing": "elasticloadbalancingv2",
	"states":               "stepfunctions",
}

func Parse(arn string) (*Arn, error) {
	if !strings.HasPrefix(arn, "arn:") {
		return nil, errors.New("invalid arn prefix")
	}
	splitArn := strings.Split(arn, ":")

	if len(splitArn) < standardLength {
		return nil, errors.New("not enough arn sections")
	}

	result := buildBaseResult(splitArn)

	if _, ok := resourceBases[result.Service]; ok {
		handleResourcePrefixARN(splitArn[resourceIndex], result)
		return result, nil
	}

	if len(splitArn) > standardLength {
		handleSevenSegmentArn(splitArn, result)
		return result, nil
	}

	resourceType, ok := ServiceBaseTypes[result.Service]
	if ok {
		handleNonStandardArn(resourceType, splitArn[resourceIndex], result)
		return result, nil
	}

	resourceId := strings.Split(strings.TrimLeft(splitArn[resourceIndex], "/"), "/")

	if len(resourceId) < 2 {
		return nil, fmt.Errorf("arn resource ID is invalid: %s", resourceId)
	}

	if len(resourceId) > 2 {
		handleRelativeResourceId(resourceId, result)
		return result, nil
	}

	// ResourceID len is 2
	handlePathStyleResourceId(resourceId, result)
	return result, nil
}

func buildBaseResult(splitArn []string) *Arn {
	service := splitArn[serviceIndex]
	if _, ok := serviceMapping[service]; ok {
		service = serviceMapping[service]
	}

	return &Arn{
		Partition: splitArn[partitionIndex],
		Service:   service,
		Region:    splitArn[regionIndex],
		AccountId: splitArn[accountIndex],
	}
}

func handleSevenSegmentArn(splitArn []string, arn *Arn) {
	arn.Resource = toSingular(splitArn[resourceIndex])
	arn.ResourceId = splitArn[resourceIdIndex]
}

func handleNonStandardArn(resourceType string, resourceId string, arn *Arn) {
	arn.Resource = resourceType
	arn.ResourceId = resourceId
}

func handlePathStyleResourceId(resourceId []string, arn *Arn) {
	arn.Resource = toSingular(resourceId[0])
	arn.ResourceId = resourceId[1]
}

func handleRelativeResourceId(resourceId []string, arn *Arn) {
	arn.Resource = toSingular(resourceId[len(resourceId)-2])
	arn.ResourceId = resourceId[len(resourceId)-1]
}

func handleResourcePrefixARN(resource string, arn *Arn) {
	split := strings.SplitN(resource, "/", 2)
	arn.Resource = split[0]
	arn.ResourceId = split[1]
}

func toSingular(word string) string {
	return strings.TrimRight(word, "s")
}
