/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/laujamie/lunchquest/internal/constants"
	"github.com/laujamie/lunchquest/internal/questrade"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// reauthenticateCmd represents the reauthenticate command
var reauthenticateCmd = &cobra.Command{
	Use:   "reauthenticate",
	Short: "Refresh Questrade API authentication tokens",
	Long: `Reauthenticate with Questrade API using the stored refresh token.

This command retrieves your saved refresh token and exchanges it for new
access and refresh tokens. Use this when your current tokens have expired
or are about to expire.

The new tokens will be stored securely in your system keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		refreshToken, err := keyring.Get(constants.SERVICE_NAME, constants.REFRESH_TOKEN_KEY)
		if err != nil {
			fmt.Println("failed to get refresh token")
			return
		}
		_, err = questrade.Authenticate(ctx, refreshToken)
		if err != nil {
			fmt.Println("failed to get new tokens")
		}
		fmt.Println("tokens updated successfully")
	},
}

func init() {
	rootCmd.AddCommand(reauthenticateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reauthenticateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reauthenticateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
