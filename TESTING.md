# Testing Real Resources

To get (potentially a lot of) resources from AWS to test against the below command can be used to pull all resources the tagging API supports. The command will write the ARNs to a file which can then be used to test parsing

`aws ec2 describe-regions --region us-west-2 | jq -r '.Regions[].RegionName' | sort  | while read REGION; do aws resourcegroupstaggingapi get-resources --region $REGION | jq --arg Region $REGION -r '.ResourceTagMappingList[]|"\(.ResourceARN)"'; done  > arns.txt`


Use the below snippet added to a test file or as a `main` to produce an errors file
that will be deduplicated based off the service:resource. Use the file to validate if the resource is available in the [Cloud Control API](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/supported-resources.html).
```go
func TestRealParsing(t *testing.T) {
	f, err := os.Open("<path-to-file>/arns.txt")
	if err != nil {
		t.FailNow()
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	errorHolder := make(map[string]struct{})
	for scanner.Scan() {
		_, err := awsresourcetypes.Lookup(scanner.Text())
		if err != nil {
			errorHolder[err.Error()] = struct{}{}
			t.Fail()
		}
	}

	file, err := os.Create("errors.txt")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for k := range errorHolder {
		fmt.Fprintln(w, k)
	}
	w.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
```
