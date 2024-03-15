package chainapi

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarcoBuarque/chain-alert/chain-alert/constants"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IChainAPI interface {
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetBlock(ctx context.Context, blockNum uint64) (*types.Block, error)
	ValidateTransaction(ctx context.Context, tx *types.Transaction) error
	WatchBlocks(blockNum chan uint64) error
}

var (
	singleton map[string]*chainAPI
)

func init() {
	singleton = make(map[string]*chainAPI, len(constants.AVALIABLE_CHAINS_ABBR))

	for _, abbr := range constants.AVALIABLE_CHAINS_ABBR {
		//"https://cloudflare-eth.com"
		client, err := ethclient.Dial(os.Getenv(fmt.Sprint("NODE_API_URL_", abbr)))
		if err != nil {
			log.Fatal(err)
		}

		//"wss://ropsten.infura.io/ws"
		websocket, err := ethclient.Dial(os.Getenv(fmt.Sprint("NODE_WEBSOCKET_URL_", abbr)))
		if err != nil {
			fmt.Println("failed to instantiate websocket", err)
		}

		singleton[abbr] = &chainAPI{apiETH: client, websocketETH: websocket}
	}
}

func API(chainAbbr string) IChainAPI {
	c, ok := singleton[chainAbbr]
	if !ok {
		return chainAPI{}
	}

	return c
}

func Close() {
	for _, c := range singleton {
		c.close()
	}
}
