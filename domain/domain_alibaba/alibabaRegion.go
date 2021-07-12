package domain_alibaba

func GetAlibabaRegionList() []string {
	return []string {
		// mainland China
		"cn-qingdao",
		"cn-beijing",
		"cn-zhangjiakou",
		"cn-huhehaote",
		"cn-wulanchabu",
		"cn-hangzhou",
		"cn-shanghai",
		"cn-shenzhen",
		"cn-heyuan",
		"cn-guangzhou",
		"cn-chengdu",
		"cn-nanjing",
		// outside mainland China
		"cn-hongkong",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-5",
		"ap-south-1",
		"ap-northeast-1",
		"us-west-1",
		"us-east-1",
		"eu-central-1",
		"eu-west-1",
		"me-east-1",
	}
}

func IsAlibabaRegionValid(region string) bool {
	regions := GetAlibabaRegionList()
	for _, r := range regions {
		if region == r {
			return true
		}
	}
	return false
}