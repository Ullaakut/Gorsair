package main

import (
	"fmt"

	"github.com/fatih/color"
)

func printSummary(targets []vulnerableDockerAPI) {
	blue := color.New(color.FgBlue, color.Underline).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	if len(targets) == 0 {
		fmt.Printf("%s No vulnerable Docker containers were found. Please make sure that your target is on an accessible network.\n", red("\xE2\x9C\x96"))
		return
	}

	for _, target := range targets {
		fmt.Printf("%s Vulnerable docker API found:\n", green("\xE2\x96\xB6"))
		fmt.Printf("    Endpoint address:\t%s\n", blue(target.Host))
		fmt.Printf("    Endpoint API port:\t%s\n", fmt.Sprint(target.Port))
		fmt.Printf("    Docker version:\t%s\n", target.DockerVersion)

		if target.SocketError != nil {
			fmt.Printf("    Docker API was unreachable:\t%s\n", red(target.SocketError))
		} else {
			fmt.Printf("    Operating system:\t%s\n", target.Info.OS)
			if len(target.Containers) > 0 {
				fmt.Printf("\n    %s running containers:\n", green(fmt.Sprint(len(target.Containers))))
				for _, container := range target.Containers {
					fmt.Printf("        %s %+v\n", container.Image, container.Mounts)
				}
			} else {
				fmt.Println("    No running containers")
			}

			if len(target.Images) > 0 {
				fmt.Printf("\n    %s available images:\n", green(len(target.Images)))
				for _, image := range target.Images {
					fmt.Printf("        %s\n", image)
				}
			} else {
				fmt.Println("    No available images")
			}

			if len(target.Containers) > 0 {
				fmt.Println("\n    To get privileged access, try running:")

				for _, container := range target.Containers {
					fmt.Printf("        docker -H %s:%d exec -it %s sh\n", target.Host, target.Port, container.ID)
				}
			}
		}
	}

	var summaryStr string
	if len(targets) == 1 {
		summaryStr = "\n%s Successful attack: %s device was accessed"
	} else {
		summaryStr = "\n%s Successful attack: %s devices were accessed"
	}

	fmt.Printf(summaryStr, green("\xE2\x9C\x94"), green(len(targets)))
}
