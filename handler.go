package main

import (
	"log"
	"os"
	"strconv"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
)

const (
	defaultRetireDecisionPeriodHour = "24" // Default: 24 hours
	defaultDryRun                   = "true"
)

// Conf -
type Conf struct {
	apiKey                   string
	retireDecisionPeriodHour int
	dryRun                   bool
}

// newConf -
func newConf() (*Conf, error) {
	apiKey, ok := os.LookupEnv("MACKEREL_API_TOKEN")
	if !ok {
		return nil, errors.Errorf("set environment variable `MACKEREL_API_TOKEN`")
	}

	r, ok := os.LookupEnv("RETIRE_DECISION_PERIOD_HOUR")
	if !ok {
		r = defaultRetireDecisionPeriodHour
	}
	ri, err := strconv.Atoi(r)
	if err != nil {
		return nil, err
	}

	d, ok := os.LookupEnv("RETIRE_DRY_RUN")
	if !ok {
		d = defaultDryRun
	}
	dr, err := strconv.ParseBool(d)
	if err != nil {
		return nil, err
	}

	return &Conf{
		apiKey: apiKey,
		retireDecisionPeriodHour: ri,
		dryRun: dr,
	}, nil
}

func handler() error {

	conf, err := newConf()
	if err != nil {
		log.Println(err)
		return err
	}

	mc := NewMackerelClient(conf.apiKey, conf.retireDecisionPeriodHour)

	hosts, err := mc.client.FindHosts(&mackerel.FindHostsParam{})
	if err != nil {
		log.Println(err)
		return err
	}
	for _, host := range hosts {
		retired, err := mc.IsRetired(host.ID)
		if err != nil {
			continue
		}
		if retired {
			if conf.dryRun {
				log.Println("[DRY RUN] Retired host:", host.ID)
			} else {
				if err := mc.client.RetireHost(host.ID); err != nil {
					log.Println(err)
					continue
				}
				log.Println("Retired host:", host.ID)
			}
		}
	}

	return nil
}
