package token

import (
	"fmt"
	"time"

	"github.com/rkonfj/peerguard/peermap/auth"
	"github.com/spf13/cobra"
)

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:   "token",
		Short: "Generate a pre-shared network secret",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterKey, err := cmd.Flags().GetString("cluster-key")
			if err != nil {
				return err
			}
			networkID, err := cmd.Flags().GetString("network")
			if err != nil {
				return err
			}
			validDuration, err := cmd.Flags().GetDuration("duration")
			if err != nil {
				return err
			}
			token, err := auth.NewAuthenticator(clusterKey).GenerateToken(networkID, validDuration)
			if err != nil {
				return err
			}
			fmt.Println(token)
			return nil
		},
	}
	Cmd.Flags().String("network", "", "network")
	Cmd.Flags().String("cluster-key", "", "key to generate token")
	Cmd.Flags().Duration("duration", 365*24*time.Hour, "secret duration to expire")

	Cmd.MarkFlagRequired("network")
	Cmd.MarkFlagRequired("cluster-key")
}
