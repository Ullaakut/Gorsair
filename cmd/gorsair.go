package main

import (
	"fmt"
	"os"

	"github.com/Ullaakut/nmap"
	"github.com/spf13/cobra"
)

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
	proxies, _ := cmd.Flags().GetStringSlice("proxies")
	decoys, _ := cmd.Flags().GetStringSlice("decoys")
	spoofIP, _ := cmd.Flags().GetString("spoofIP")
	spoofMAC, _ := cmd.Flags().GetString("spoofMAC")
	iface, _ := cmd.Flags().GetString("interface")
	verbose, _ := cmd.Flags().GetBool("verbose")

	options := []func(*nmap.Scanner){
		nmap.WithTargets(targets...),
		nmap.WithPorts(ports...),
		nmap.WithTimingTemplate(nmap.Timing(speed)),
		nmap.WithSYNScan(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithServiceInfo(),
	}

	if len(decoys) != 0 {
		options = append(options, nmap.WithDecoys(decoys...))
	}

	if len(proxies) != 0 {
		options = append(options, nmap.WithProxies(proxies...))
	}

	if spoofIP != "" {
		options = append(options, nmap.WithSpoofIPAddress(spoofIP))
	}

	if spoofMAC != "" {
		options = append(options, nmap.WithSpoofMAC(spoofMAC))
	}

	if iface != "" {
		options = append(options, nmap.WithInterface(iface))
	}

	scanner, err := nmap.NewScanner(options...)
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

				api := vulnerableDockerAPI{
					Endpoint:      fmt.Sprint("tcp://", addr.String(), ":", port.ID),
					Host:          addr.String(),
					Port:          port.ID,
					DockerVersion: version,
				}

				err := gatherInformation(&api)
				if err != nil {
					if verbose {
						fmt.Printf("Unable to access the docker API through endpoint %q: %v", api.Endpoint, err)
					}

					api.SocketError = err
				}

				vulnerableTargets = append(vulnerableTargets, api)
			}
		}
	}

	printSummary(vulnerableTargets)
}
