package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show domain.",
	Long:  `After execution, it will return to the existing data.`,
	Run: func(cmd *cobra.Command, args []string) {
		domains := viper.GetStringSlice(DOMAINS)
		for _, arg := range args {
			for i, domain := range domains {
				if strings.EqualFold(domain, arg) {
					domains = append(domains[:i], domains[i:]...)
				}
			}
		}
		viper.Set(DOMAINS, domains)
		fmt.Println("The modified domain name list is: ")
		for _, domain := range domains {
			fmt.Println(domain)
		}
	},
}
