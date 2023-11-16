package parser

import (
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

func Parse(arn string) Arn {
	splitArn := strings.Split(arn, ":")

	result := buildBaseResult(splitArn)

	if len(splitArn) > standardLength {
		handleSevenSegmentArn(splitArn, result)
		return *result
	}

	resourceType, ok := ServiceBaseTypes[splitArn[serviceIndex]]

	if ok {
		handleNonStandardArn(resourceType, splitArn[resourceIndex], result)
		return *result
	}

	resourceId := strings.Split(strings.TrimLeft(splitArn[resourceIndex], "/"), "/")

	if len(resourceId) == 2 {
		handlePathStyleResourceId(resourceId, result)
		return *result
	}

	if len(resourceId) > 2 {
		handleRelativeResourceId(resourceId, result)
		return *result
	}

	return Arn{}
}

func buildBaseResult(splitArn []string) *Arn {
	return &Arn{
		Partition: splitArn[partitionIndex],
		Service:   splitArn[serviceIndex],
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

func toSingular(word string) string {
	return strings.TrimRight(word, "s")
}
