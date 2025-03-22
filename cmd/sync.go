/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/laujamie/lunchquest/internal/constants"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := keyring.Get(constants.SERVICE_NAME, constants.ACCESS_TOKEN_KEY)
		if err != nil {
			fmt.Println("Error getting access token:", err)
			return
		}
		fmt.Println(accessToken)
		req, err := http.NewRequest(http.MethodGet, "https://api03.iq.questrade.com/v1/accounts", nil)
		if err != nil {
			fmt.Println("failed to create req")
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("failed to execute req")
			return
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read body")
			return
		}

		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
