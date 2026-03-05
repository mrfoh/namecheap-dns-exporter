package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrfoh/namecheap-dns-exporter/internal/namecheap"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.0"

var (
	domain string
	format string
)

var rootCmd = &cobra.Command{
	Use:     "namecheap-dns-exporter <input-file>",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"ncdns-exporter", "ncdns"},
	Version: VERSION,
	Short:   "A CLI tool to convert a Namecheap DNS JSON export into a zone file.",
	RunE:    runExporter,
}

func init() {
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "domain name for the zone file (required)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "zone", "output format: zone or route53")
	rootCmd.MarkFlagRequired("domain")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}

func runExporter(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	f, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	return namecheap.Export(f, os.Stdout, domain, namecheap.Format(format))
}
