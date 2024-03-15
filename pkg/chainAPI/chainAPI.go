package chainapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ETH BASED
type chainAPI struct {
	apiETH       *ethclient.Client
	websocketETH *ethclient.Client
}

func (controller chainAPI) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	blockNum, err := controller.apiETH.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blockNum, nil
}
func (controller chainAPI) GetBlock(ctx context.Context, blockNum uint64) (*types.Block, error) {
	blockNumber := big.NewInt(int64(blockNum))

	block, err := controller.apiETH.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return &types.Block{}, err
	}

	return block, nil
}

func (controller chainAPI) ValidateTransaction(ctx context.Context, tx *types.Transaction) error {
	receipt, err := controller.apiETH.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}

	if receipt.Status == 0 {
		return ErrTxFailed
	}

	return nil
}

func (controller chainAPI) WatchBlocks(blockNum chan uint64) error {
	if controller.websocketETH == nil {
		return fmt.Errorf("websocket not supported")
	}

	channel := make(chan *types.Header)
	sub, err := controller.websocketETH.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	for {
		select {
		case block := <-channel:
			blockNum <- block.Number.Uint64()

		case err := <-sub.Err():
			return err
		}
	}
}

func (controller chainAPI) close() { controller.apiETH.Close() }
