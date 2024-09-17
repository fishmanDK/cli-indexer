package main

import "github.com/fishmanDK/cli-indexer/cmd"

func main() {
	cmd.Execute()
}



// package main

// import (
// 	"log"

// 	"github.com/fishmanDK/cli-indexer/indexer"
// 	"github.com/spf13/pflag"
// 	"github.com/spf13/viper"
// )

// func main() {
// 	rpcUrl := pflag.String("rpc", "", "url blockchain rpc")
// 	startBlock := pflag.Int64("start", 0, "number of block where parsing begins")
// 	outputFile := pflag.String("out", "", "name of file")

// 	pflag.Parse()
// 	viper.BindPFlags(pflag.CommandLine)

// 	if *rpcUrl == "" || *startBlock <= 0 || *outputFile == ""{
// 		log.Fatal("all flags must be specified")
// 	}
	
// 	err := indexer.RunIndexer(*rpcUrl, *startBlock, *outputFile)
// 	if err != nil{
// 		log.Fatalf("error indexer: %v", err)
// 	}
// }

