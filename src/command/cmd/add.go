package cmd

import (
	"command/service/windows/host"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	rootCmd.AddCommand(addCmd)
	//rootCmd.AddCommand(addAmazonDomainCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add domain name to list.",
	Long:  `Duplicates will be ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := viper.GetStringSlice(DOMAINS)
		for _, arg := range args {
			if !StringArrayContains(arg, domains) {
				domains = append(domains, arg)
				host.AutomaticModification(arg)
			}
		}
		viper.Set(DOMAINS, domains)
		fmt.Println("The modified domain name list is: ")
		for _, domain := range domains {
			fmt.Println(domain)
		}
	},
}

var addAmazonDomainCmd = &cobra.Command{
	Use:   "add ",
	Short: "Add domain name to list.",
	Long:  `Duplicates will be ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := viper.GetStringSlice(DOMAINS)
		for _, arg := range args {
			if !StringArrayContains(arg, domains) {
				domains = append(domains, arg)
			}
		}
		viper.Set(DOMAINS, domains)
		fmt.Println("The modified domain name list is: ")
		for _, domain := range domains {
			fmt.Println(domain)
		}
	},
}

func StringArrayContains(s string, arr []string) bool {
	for _, s2 := range arr {
		if strings.EqualFold(s, s2) {
			return true
		}
	}
	return false
}
