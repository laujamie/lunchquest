/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/laujamie/lunchquest/internal/questrade"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command used to authenticate with the Questrade API.
// It exchanges a refresh token for an access token and stores the authentication credentials
// for subsequent API calls. The refresh token can be provided using the --refresh_token
// flag. Upon successful authentication, a confirmation message is displayed.
//
// The command will timeout after 30 seconds if the authentication process takes too long.
// If authentication fails, an error message will be displayed indicating what went wrong.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Questrade",
	Long: `Login to the Questrade API using your refresh token.

This command authenticates your session with Questrade by exchanging your 
refresh token for an access token. Authentication is required before you 
can make any API calls to access your account data or market information.

Example usage:
  lunchquest login -r YOUR_REFRESH_TOKEN

You can obtain a refresh token from the Questrade App Hub after registering
your application. The token is valid for a limited time and will need to be
refreshed periodically.`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshToken, err := cmd.Flags().GetString("refresh_token")
		if err != nil {
			fmt.Println("Error accessing refresh token:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
		defer cancel()

		_, err = questrade.Authenticate(ctx, refreshToken)
		if err != nil {
			fmt.Println("Error during authentication:", err)
			return
		}
		fmt.Println("Successfully authenticated")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringP("refresh_token", "r", "", "Refresh token for authentication")
	loginCmd.MarkFlagRequired("refresh_token")
}
