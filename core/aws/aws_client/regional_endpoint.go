package aws_client

import (
	"fmt"
	"leapp_daemon/custom_error"
)

var regionalEndpoints = map[string]string {
	"af-south-1": "https://sts.af-south-1.amazonaws.com",
	"ap-east-1": "https://sts.ap-east-1.amazonaws.com",
	"ap-northeast-1": "https://sts.ap-northeast-1.amazonaws.com",
	"ap-northeast-2": "https://sts.ap-northeast-2.amazonaws.com",
	"ap-northeast-3": "https://sts.ap-northeast-3.amazonaws.com",
	"ap-south-1": "https://sts.ap-south-1.amazonaws.com",
	"ap-southeast-1": "https://sts.ap-southeast-1.amazonaws.com",
	"ap-southeast-2": "https://sts.ap-southeast-2.amazonaws.com",
	"ca-central-1": "https://sts.ca-central-1.amazonaws.com",
	"cn-north-1": "https://sts.cn-north-1.amazonaws.com.cn",
	"cn-northwest-1": "https://sts.cn-northwest-1.amazonaws.com.cn",
	"eu-central-1": "https://sts.eu-central-1.amazonaws.com",
	"eu-north-1": "https://sts.eu-north-1.amazonaws.com",
	"eu-south-1": "https://sts.eu-south-1.amazonaws.com",
	"eu-west-1": "https://sts.eu-west-1.amazonaws.com",
	"eu-west-2": "https://sts.eu-west-2.amazonaws.com",
	"eu-west-3": "https://sts.eu-west-3.amazonaws.com",
	"me-south-1": "https://sts.me-south-1.amazonaws.com",
	"sa-east-1": "https://sts.sa-east-1.amazonaws.com",
	"us-east-1": "https://sts.us-east-1.amazonaws.com",
	"us-east-2": "https://sts.us-east-2.amazonaws.com",
	"us-gov-east-1": "https://sts.us-gov-east-1.amazonaws.com",
	"us-gov-west-1": "https://sts.us-gov-west-1.amazonaws.com",
	"us-west-1": "https://sts.us-west-1.amazonaws.com",
	"us-west-2": "https://sts.us-west-2.amazonaws.com",
}

func GetRegionalEndpoint(region *string) (*string, error) {
	isRegionValid := IsRegionValid(*region)
	if !isRegionValid {
		return nil, custom_error.NewUnprocessableEntityError(fmt.Errorf("Region " + *region + " not valid"))
	}

	i, ok := regionalEndpoints[*region]

	if ok {
		return &i, nil
	} else {
		return &i, custom_error.NewNotFoundError(fmt.Errorf("No such endpoint for region " + *region))
	}
}
