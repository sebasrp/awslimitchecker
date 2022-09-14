package services

type AWSQuotaInfo struct {
	Service    string  // service the quota applies to
	Name       string  // the name of the aws service resource the usage is for
	ResourceId string  // if there can be multiple usages for one quota, aws id (Cloudformation format)
	Region     string  // the region this quota applies to
	Quotacode  string  // servicequota code
	QuotaValue float64 // the quota value
	UsageValue float64 // the usage value
	Unit       string  // unit of the quota/usage
	Global     bool    // whether the quota is global or not
}

type Svcquota interface {
	// Get Usage retrieve the quotas and usage for the given service
	GetUsage() []AWSQuotaInfo

	// GetAllDefaultQuotas retrieves all the default quotas for the given
	// service. Usage of those resources are not retrieved/calculated
	GetAllDefaultQuotas() map[string]AWSQuotaInfo

	// GetRequiredPermissions returns a list of the IAM permissions required
	// to retrieve the usage for this service.
	GetRequiredPermissions() []string
}
