package main

import (
	"fmt"

	"github.com/Ullaakut/disgo"
	"github.com/Ullaakut/disgo/style"
)

func printSummary(targets []vulnerableDockerAPI) {
	term := disgo.NewTerminal()

	if len(targets) == 0 {
		term.Errorln(style.Failure(style.SymbolCross), "No vulnerable Docker containers were found. Please make sure that your target is on an accessible network.")
		return
	}

	for _, target := range targets {
		term.Infoln(style.Success(style.SymbolRightTriangle), "Vulnerable docker API found:")
		term.Infof("    Endpoint address:\t%s\n", style.Link(target.Host))
		term.Infof("    Endpoint API port:\t%v\n", target.Port)
		term.Infof("    Docker version:\t%s\n", target.DockerVersion)

		if target.SocketError != nil {
			term.Infof("    Docker API was unreachable:\t%s\n", style.Failure(target.SocketError))
		} else {
			term.Infof("    Operating system:\t%s\n", target.Info.OS)
			if len(target.Containers) > 0 {
				term.Infof("\n    %s running containers:\n", style.Success(len(target.Containers)))
				for _, container := range target.Containers {
					term.Infof("        %s %+v\n", container.Image, container.Ports)
				}
			} else {
				fmt.Println("    No running containers")
			}

			if len(target.Images) > 0 {
				term.Infof("\n    %s available images:\n", style.Success(len(target.Images)))
				for _, image := range target.Images {
					term.Infof("        %s\n", image)
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

	term.Infof(summaryStr, style.Success(style.SymbolCheck), style.Success(len(targets)))
}
