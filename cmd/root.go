package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fishmanDK/cli-indexer/indexer"
	"github.com/spf13/cobra"
)

var rpcURL string
var startBlock int64
var outputFile string

var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "Indexer CLI application",
	Long:  "A simple CLI application for indexing blocks.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Indexer application running...")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the indexer",
	Run: func(cmd *cobra.Command, args []string) {
		if rpcURL == "" || startBlock < 0 || outputFile == ""{
			log.Fatal("all flags must be specified")
		}

		if err := indexer.RunIndexer(rpcURL, startBlock, outputFile); err != nil {
			fmt.Printf("Error running indexer: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&rpcURL, "rpc", "", "RPC URL for Ethereum client")
	runCmd.Flags().Int64Var(&startBlock, "start", 0, "Starting block number")
	runCmd.Flags().StringVar(&outputFile, "out", "./logs/blocks.log", "Output file for blocks")
}

func Execute() error {
	return rootCmd.Execute()
}