/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"time"

	"github.com/laujamie/lunchquest/internal/questrade"
	"github.com/spf13/cobra"
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
		defer cancel()

		oauthToken, err := questrade.GetStoredAuthToken()
		if err != nil {
			log.Fatalf("failed to get refresh token: %v", err)
			return
		}

		_, err = questrade.Authenticate(ctx, oauthToken.RefreshToken)
		if err != nil {
			log.Fatalf("failed to get new tokens: %v", err)
			return
		}
		log.Println("tokens updated successfully")
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
