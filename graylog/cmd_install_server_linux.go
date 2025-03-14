package graylog

import (
	"github.com/spf13/cobra"
)

var (
	tlsPublicChainPath string
	tlsPrivateKeyPath  string
)

var installServerCmd = &cobra.Command{
	Use:   "install-server",
	Short: "Install the Graylog server",
	Run: func(cmd *cobra.Command, args []string) {
		installServer(tlsPublicChainPath, tlsPrivateKeyPath)
	},
}

func setupInstallServerCmd(cmd *cobra.Command) {
	installServerCmd.Flags().StringVar(&tlsPublicChainPath, "tls-public-chain", "graylog.internal.bundle.crt", "Path to TLS public certificate chain")
	installServerCmd.Flags().StringVar(&tlsPrivateKeyPath, "tls-private-key", "graylog.internal.key", "Path to TLS private key")
	installServerCmd.MarkFlagRequired("tls-public-chain")
	installServerCmd.MarkFlagRequired("tls-private-key")

	cmd.AddCommand(installServerCmd)
}
