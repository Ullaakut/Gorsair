package main

import (
	"fmt"

	"github.com/Ullaakut/nmap"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gorsair",
	Short: "Gorsair hacks into containers that expose their Docker APIs",
	Long: `Gorsair discovers and hacks into vulnerable containers that expose their Docker APIs.

	Do not use this software on a network that you do not own.`,
	Run: gorsair,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkRequiredFlags(cmd.Flags())
	},
}

func checkRequiredFlags(flags *pflag.FlagSet) error {
	targets, _ := flags.GetStringSlice("targets")
	if len(targets) == 0 {
		return fmt.Errorf("'targets' argument is required")
	}

	speed, _ := flags.GetInt("speed")
	if nmap.Timing(speed) < nmap.TimingSneaky || nmap.Timing(speed) > nmap.TimingFastest {
		return fmt.Errorf("speed %d is invalid: value should be between %d and %d", speed, nmap.TimingSneaky, nmap.TimingAggressive)
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().StringSliceP("targets", "t", []string{}, "List of targets to scan in nmap format (see https://nmap.org/book/man-target-specification.html)")
	rootCmd.PersistentFlags().StringSliceP("ports", "p", []string{"2375", "2376"}, "List of ports to scan")
	rootCmd.PersistentFlags().IntP("speed", "s", 4, "Speed at which to scan the network. Lower is stealthier (see https://nmap.org/book/man-performance.html)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")

	viper.BindPFlag("targets", rootCmd.PersistentFlags().Lookup("targets"))
	viper.BindPFlag("ports", rootCmd.PersistentFlags().Lookup("ports"))
	viper.BindPFlag("speed", rootCmd.PersistentFlags().Lookup("speed"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func execute() error {
	return rootCmd.Execute()
}
