package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mackerelio/mackerel-client-go"
)

func main() {
	apiKey, ok := os.LookupEnv("MACKEREL_API_TOKEN")
	if !ok {
		log.Fatal("set environment variable `MACKEREL_API_TOKEN`")
	}

	r, ok := os.LookupEnv("RETIRE_DECISION_PERIOD_HOUR")
	if !ok {
		r = defaultRetireDecisionPeriodHour
	}
	retireDecisionPeriodHour, err := strconv.Atoi(r)
	if err != nil {
		log.Fatal(err)
	}

	d, ok := os.LookupEnv("RETIRE_DRY_RUN")
	if !ok {
		d = "true"
	}
	dryRun, err := strconv.ParseBool(d)
	if err != nil {
		log.Fatal(err)
	}

	mc := NewMackerelClient(apiKey, retireDecisionPeriodHour)

	hosts, err := mc.client.FindHosts(&mackerel.FindHostsParam{})
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {

		retired, err := mc.IsRetired(host.ID)
		if err != nil {
			continue
		}
		if retired {
			fmt.Println("Retired host:", host.ID)
			if !dryRun {
				fmt.Println(dryRun)
				//if err := mc.client.RetireHost(host.ID); err != nil {
				//	fmt.Println(err)
				//	continue
				//}
			}
		}
	}

}
