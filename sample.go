
package main

import (
    "context"
    "fmt"
    ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
    "github.com/HydroProtocol/ethereum-watcher/blockchain"
    "github.com/HydroProtocol/ethereum-watcher/plugin"
    "github.com/shopspring/decimal"
    "github.com/sirupsen/logrus"

    //"github.com/HydroProtocol/ethereum-watcher/structs"
)


func main() {
    api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
    w := ethereum_watcher.NewHttpBasedEthWatcher(context.Background(), api)

    // we use TxReceiptPlugin here
    w.RegisterTxReceiptPlugin(plugin.NewERC20TransferPlugin(
        func(token, from, to string, amount decimal.Decimal, isRemove bool) {

            logrus.Infof("New ERC20 Transfer >> token(%s), %s -> %s, amount: %s, isRemoved: %t",
                token, from, to, amount, isRemove)

        },
    ))

    // we use BlockPlugin here
    w.RegisterBlockPlugin(plugin.NewBlockNumPlugin(func(i uint64, b bool) {
        fmt.Println(">>", i, b)
    }))


    usdtContractAdx := "0xdac17f958d2ee523a2206206994597c13d831ec7"

    // ERC20 Transfer Event
    topicsInterestedIn := []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}

    handler := func(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {
        logrus.Infof("USDT Transfer count: %d, %d -> %d", len(receiptLogs), from, to)
        return nil
    }

    // query for USDT Transfer Events
    receiptLogWatcher := ethereum_watcher.NewReceiptLogWatcher(
        context.TODO(),
        api,
        -1,
        usdtContractAdx,
        topicsInterestedIn,
        handler,
        ethereum_watcher.ReceiptLogWatcherConfig{
            StepSizeForBigLag:               5,
            IntervalForPollingNewBlockInSec: 5,
            RPCMaxRetry:                     3,
            ReturnForBlockWithNoReceiptLog:  true,
        },
    )

    receiptLogWatcher.Run()


    w.RunTillExit()
}
