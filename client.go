package main

import (
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

const (
	ec2CpuUsageMetricName         = "custom.ec2.cpu.used"
	gceCpuUsageMetricName         = "cpu.user.percentage" // Since there is no GCE specific metrics, I use normal CPU metrics.
	rdsCpuUsageMetricName         = "custom.rds.cpu.used"
	redShiftCpuUsageMetricName    = "custom.redshift.cpu.used"
	elastiCacheCpuUsageMetricName = "custom.elasticache.cpu.used"
)

// MackerelClient -
type MackerelClient struct {
	from   int64
	to     int64
	client *mackerel.Client
}

// NewMackerelClient -
func NewMackerelClient(apiKey string, retireDecisionPeriodHour int) *MackerelClient {
	now := time.Now()
	return &MackerelClient{
		to:     now.Unix(),
		from:   now.Add(time.Duration(-1*retireDecisionPeriodHour) * time.Hour).Unix(),
		client: mackerel.NewClient(apiKey),
	}
}

// IsRetired returns true if the host is a target of retirement.
// Hosts that do not have metrics for the specified period are subject to retirement.
func (mc *MackerelClient) IsRetired(hostId string) (bool, error) {
	h, err := mc.client.FindHost(hostId)
	if err != nil {
		return false, err
	}
	if h.Meta.Cloud == nil {
		return false, nil
	}

	var met string
	switch h.Meta.Cloud.Provider {
	case "ec2":
		met = ec2CpuUsageMetricName
	case "gce":
		met = gceCpuUsageMetricName
	case "rds":
		met = rdsCpuUsageMetricName
	case "redshift":
		met = redShiftCpuUsageMetricName
	case "elasticache":
		met = elastiCacheCpuUsageMetricName
	default:
		return false, nil
	}

	// If there are no metrics for the specified period, retire the host
	values, err := mc.client.FetchHostMetricValues(hostId, met, mc.from, mc.to)
	if err != nil {
		return false, err
	}
	if len(values) == 0 {
		return true, nil
	}

	return false, nil
}
