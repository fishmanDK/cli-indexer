// package indexer

// import (
// 	"context"
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func RunIndexer(rpcURL string, startBlock int64, outputFile string) error {
// 	client, err := ethclient.Dial(rpcURL)
// 	if err != nil {
// 		return fmt.Errorf("RPC connection error: %v", err)
// 	}

// 	var fileBlocksLog *os.File
// 	_, err = os.Stat(outputFile)
// 	if err != nil{
// 		if os.IsNotExist(err) {
// 			fmt.Printf("no such file: %v\n", err)
// 			fileBlocksLog, err = os.Create("./logs/blocks.log")
// 			if err != nil {
// 				return fmt.Errorf("failed to create file: %v", err)
// 			}
// 		} else {
// 			return err
// 		}
// 	} else {
// 		fileBlocksLog, err = os.Open(outputFile)
// 		if err != nil {
// 			return fmt.Errorf("failed open file: %v", err)
// 		}
// 	}

// 	defer fileBlocksLog.Close()

// 	fileErrorsLog, err := os.Create("./logs/errors.log")
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %s", "errors.log")
// 	}
// 	defer fileErrorsLog.Close()

// 	lastBlock, err := client.BlockNumber(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	cnt := 0

// 	stop := make(chan struct{})
// 	go func(stop chan struct{}, cnt *int) {
// 		startTime := time.Now()
// 		ticker := time.NewTicker(100 * time.Millisecond)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-stop:
// 				return
// 			case <-ticker.C:
// 				currentTime := time.Since(startTime)

// 				hours := int(currentTime.Hours())
// 				minutes := int(currentTime.Minutes()) % 60
// 				seconds := int(currentTime.Seconds()) % 60
// 				milliseconds := int((currentTime.Milliseconds() % 1000) / 10)

// 				// fmt.Fprintf(os.Stdout, "\033[2K\r")

// 				// fmt.Fprintf(os.Stdout, "Time: %02d:%02d:%02d.%02d", hours, minutes, seconds, milliseconds)

// 				fmt.Fprintf(os.Stdout, "\033[2K\r") // Очистка строки
// 				fmt.Fprintf(os.Stdout, "Time: %02d:%02d:%02d.%02d | Errors: %d", hours, minutes, seconds, milliseconds, *cnt) // Таймер и количество ошибок
// 				os.Stdout.Sync()
// 			}
// 		}
// 	}(stop, &cnt)

// 	ch := make(chan int64, 100)

// 	wg := &sync.WaitGroup{}
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func(ch chan int64) {
// 			defer wg.Done()
// 			for numOfBlock := range ch {
// 				headerBlock, err := client.HeaderByNumber(context.Background(), big.NewInt(numOfBlock))
// 				if err != nil {
// 					for attempt := 0; attempt < maxRetries; attempt++ {
// 						headerBlock, err = client.HeaderByNumber(context.Background(), big.NewInt(numOfBlock))
// 						if err == nil {
// 							break
// 						}
// 						time.Sleep(retryDelay)
// 					}

// 					if headerBlock == nil {
// 						cnt++
// 						fileErrorsLog.WriteString(fmt.Sprintf("Failed to get info about block %d:\n", numOfBlock))
// 						continue
// 					}
// 				}

// 				blockInfo := fmt.Sprintf("Number: %d Hash: %s TxCount: %d TimeStamp: %d\n", numOfBlock, headerBlock.Hash().Hex(), headerBlock.Number.Uint64(), headerBlock.Time)
// 				fileBlocksLog.WriteString(blockInfo)

// 				time.Sleep(parceBlockSleep)
// 			}
// 		}(ch)
// 	}

// 	for numBlock := startBlock; numBlock < int64(lastBlock); numBlock++ {
// 		time.Sleep(time.Microsecond * 200)
// 		ch <- numBlock
// 	}

// 	close(ch)
// 	wg.Wait()

// 	close(stop)

// 	fmt.Printf("\nTotal: %d\nSuccess: %d\nFailed: %d\n", lastBlock-uint64(startBlock), (lastBlock-uint64(startBlock))-uint64(cnt), cnt)
// 	fmt.Println("Last block: %d", lastBlock)
// 	return nil
// }

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

	var fileBlocksLog *os.File
	_, err = os.Stat(outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("no such file: %v\n", err)
			fileBlocksLog, err = os.Create("./logs/blocks.log")
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
		} else {
			return err
		}
	} else {
		fileBlocksLog, err = os.Open(outputFile)
		if err != nil {
			return fmt.Errorf("failed open file: %v", err)
		}
	}
	defer fileBlocksLog.Close()
	fileErrorsLog, err := os.Create("./logs/errors.log")
	if err != nil {
		return fmt.Errorf("failed to create file: %s", "errors.log")
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
