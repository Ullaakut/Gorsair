package main

import (
	"fmt"
	"os"

	"github.com/Ullaakut/nmap"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type vulnerableDockerAPI struct {
	Host          string
	DockerVersion string
	Port          uint16
}

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func gorsair(cmd *cobra.Command, args []string) {
	targets, _ := cmd.Flags().GetStringSlice("targets")
	ports, _ := cmd.Flags().GetStringSlice("ports")
	speed, _ := cmd.Flags().GetInt("speed")
	verbose, _ := cmd.Flags().GetBool("verbose")

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targets...),
		nmap.WithPorts(ports...),
		nmap.WithTimingTemplate(nmap.Timing(speed)),
		nmap.WithSYNScan(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithServiceInfo(),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := startSpinner(verbose)

	updateSpinner(w, "Scanning targets...", verbose)

	results, err := scanner.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clearOutput(w, verbose)

	if verbose {
		fmt.Printf("Scan successful. Hosts: %+v\n", results.Stats.Hosts)
	}

	var vulnerableTargets []vulnerableDockerAPI
	for _, host := range results.Hosts {
		for _, addr := range host.Addresses {
			for _, port := range host.Ports {
				if verbose {
					fmt.Printf("Port %d on host %s: %s\n", port.ID, addr, port.Status())
				}

				if port.Status() != nmap.Open {
					continue
				}

				if port.Service.Name != "docker" {
					continue
				}

				version := port.Service.Version
				if version == "" {
					version = "UNKNOWN"
				}

				vulnerableTargets = append(vulnerableTargets, vulnerableDockerAPI{
					Host:          addr.String(),
					Port:          port.ID,
					DockerVersion: version,
				})
			}
		}
	}

	printSummary(vulnerableTargets)
}

func printSummary(targets []vulnerableDockerAPI) {
	blue := color.New(color.FgBlue, color.Underline).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	if len(targets) == 0 {
		fmt.Printf("%s No vulnerable Docker containers were found. Please make sure that your target is on an accessible network.\n", red("\xE2\x9C\x96"))
		return
	}

	for _, container := range targets {
		fmt.Printf("%s Vulnerable container found:\n", green("\xE2\x96\xB6"))
		fmt.Printf("    Container address:\t%s\n", blue(container.Host))
		fmt.Printf("    Container API port:\t%s\n", fmt.Sprint(container.Port))
		fmt.Printf("    Docker version:\t%s\n", container.DockerVersion)

		// Docker commands
		fmt.Printf("\n    You can get more information from this vulnerable container by running the following commands:\n")
		fmt.Printf("        docker -H %s:%d info\t\tWill give you system-wide information.\n", container.Host, container.Port)
		fmt.Printf("        docker -H %s:%d ps -a\t\tWill list all of the containers running and stopped.\n", container.Host, container.Port)
		fmt.Printf("        docker -H %s:%d images\t\tWill list all of the images that are available on this docker machine.\n", container.Host, container.Port)
	}

	var summaryStr string
	if len(targets) == 1 {
		summaryStr = "\n%s Successful attack: %s device was accessed"
	} else {
		summaryStr = "\n%s Successful attack: %s devices were accessed"
	}

	fmt.Printf(summaryStr, green("\xE2\x9C\x94"), green(len(targets)))
}
