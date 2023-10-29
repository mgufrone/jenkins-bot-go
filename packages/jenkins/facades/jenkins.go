package facades

import (
	"log"

	"mgufrone.dev/job-tracking/packages/jenkins"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
)

func Jenkins() contracts.Jenkins {
	instance, err := jenkins.App.Make(jenkins.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Jenkins)
}

func Handler() contracts.Handler {
	instance, err := jenkins.App.Make(jenkins.BindingHandler)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Handler)
}
