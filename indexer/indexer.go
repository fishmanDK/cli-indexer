package indexer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func RunIndexer(rpcURL string, startBlock int64, outputFile string) error {
	client, err := connectToRPC(rpcURL)
	if err != nil {
		return err
	}

	fileBlocksLog, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer fileBlocksLog.Close()
	fileErrorsLog, err := os.Open("./logs/errors.log")
	if err != nil {
		return fmt.Errorf("failed to open file errors.log: %v", err)
	}
	defer fileErrorsLog.Close()

	lastBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	cnt := 0

	stopChan := make(chan struct{})
	go displayTimeAndErrors(stopChan, &cnt)

	ch := make(chan int64, 100)

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(ch chan int64) {
			defer wg.Done()
			processBlocks(client, ch, stopChan, fileBlocksLog, &cnt)
		}(ch)
	}

	for numBlock := startBlock; numBlock < int64(lastBlock); numBlock++ {
		time.Sleep(time.Microsecond * 200)
		ch <- numBlock
	}

	close(ch)
	wg.Wait()

	close(stopChan)

	fmt.Printf("\nTotal: %d\nSuccess: %d\nFailed: %d\n", lastBlock-uint64(startBlock), (lastBlock-uint64(startBlock))-uint64(cnt), cnt)
	fmt.Printf("Last block: %d\n", lastBlock)
	return nil
}

func connectToRPC(rpcURL string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("RPC connection error: %v", err)
	}
	return client, nil
}
