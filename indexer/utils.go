package indexer

import (
    "context"
    "fmt"
    "math/big"
    "os"
    "time"

    "github.com/ethereum/go-ethereum/ethclient"
)

func openLogFile(filename string) (*os.File, error) {
    _, err := os.Stat(filename)
    if err != nil && os.IsNotExist(err) {
        file, err := os.Create(filename)
        if err != nil {
            return nil, fmt.Errorf("failed to create log file: %v", err)
        }
        return file, nil
    }

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %v", err)
    }
    return file, nil
}

func displayTimeAndErrors(stopChan chan struct{}, cnt *int) {
    startTime := time.Now()
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-stopChan:
            return
        case <-ticker.C:
            currentTime := time.Since(startTime)

            hours := int(currentTime.Hours())
            minutes := int(currentTime.Minutes()) % 60
            seconds := int(currentTime.Seconds()) % 60
            milliseconds := int((currentTime.Milliseconds() % 1000) / 10)

            fmt.Fprintf(os.Stdout, "\033[2K\r")
            fmt.Fprintf(os.Stdout, "Time: %02d:%02d:%02d.%02d | Errors: %d", hours, minutes, seconds, milliseconds, *cnt)
            os.Stdout.Sync()
        }
    }
}

func processBlocks(client *ethclient.Client, ch chan int64, stopChan <-chan struct{}, fileBlocksLog *os.File, cnt *int) {
    for blockNum := range ch {
        select {
        case <-stopChan:
            return
        default:
            headerBlock, err := client.HeaderByNumber(context.Background(), big.NewInt(blockNum))
            if err != nil || headerBlock == nil {
                for attempt := 0; attempt < maxRetries; attempt++ {
                    headerBlock, err = client.HeaderByNumber(context.Background(), big.NewInt(blockNum))
                    if headerBlock != nil && err == nil {
                        break
                    }
                    time.Sleep(retryDelay)
                }

                if headerBlock == nil { 
                    (*cnt)++
                    continue 
                }
            }

            blockInfo := fmt.Sprintf("Number: %d Hash: %s TxCount: %d TimeStamp: %d\n", blockNum, headerBlock.Hash().Hex(), headerBlock.Number.Uint64(), headerBlock.Time)
            fileBlocksLog.WriteString(blockInfo)

            time.Sleep(parceBlockSleep) 
        }
    }
}