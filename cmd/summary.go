package main

import (
	"fmt"
	"os"

	"github.com/ullaakut/disgo/logger"
	"github.com/ullaakut/disgo/symbol"
)

func printSummary(targets []vulnerableDockerAPI) {
	log, err := logger.New(os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create logger:", err)
	}

	if len(targets) == 0 {
		log.Errorln(logger.Failure(symbol.Cross), "No vulnerable Docker containers were found. Please make sure that your target is on an accessible network.")
		return
	}

	for _, target := range targets {
		log.Infoln(logger.Success(symbol.RightTriangle), "Vulnerable docker API found:")
		log.Infof("    Endpoint address:\t%s\n", logger.Link(target.Host))
		log.Infof("    Endpoint API port:\t%v\n", target.Port)
		log.Infof("    Docker version:\t%s\n", target.DockerVersion)

		if target.SocketError != nil {
			log.Infof("    Docker API was unreachable:\t%s\n", logger.Failure(target.SocketError))
		} else {
			log.Infof("    Operating system:\t%s\n", target.Info.OS)
			if len(target.Containers) > 0 {
				log.Infof("\n    %s running containers:\n", logger.Success(len(target.Containers)))
				for _, container := range target.Containers {
					log.Infof("        %s %+v\n", container.Image, container.Ports)
				}
			} else {
				fmt.Println("    No running containers")
			}

			if len(target.Images) > 0 {
				log.Infof("\n    %s available images:\n", logger.Success(len(target.Images)))
				for _, image := range target.Images {
					log.Infof("        %s\n", image)
				}
			} else {
				fmt.Println("    No available images")
			}
		}
	}

	var summaryStr string
	if len(targets) == 1 {
		summaryStr = "\n%s Successful attack: %s device was accessed"
	} else {
		summaryStr = "\n%s Successful attack: %s devices were accessed"
	}

	log.Infof(summaryStr, logger.Success(symbol.Check), logger.Success(len(targets)))
}
