package listener

import (
	"context"
	"errors"
	"fmt"

	chainapi "github.com/MarcoBuarque/chain-alert/chain-alert/pkg/chainAPI"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

func processBlock(abbr string) {
	ctx := context.Background()
	lastSavedBlock := uint64(19440134) // TODO: GET DATA FROM CACHE

	lastBlock, err := chainapi.API(abbr).GetLastBlockNumber(ctx)
	if err != nil {
		return
	}

	txs, err := getTransactions(ctx, abbr, lastBlock, lastSavedBlock)
	if err != nil {
		return
	}

	txs, err = filterTransactions(txs)
	if err != nil {
		return
	}

	if err := processTransactions(ctx, abbr, txs); err != nil {
		fmt.Println("failed when try to process transactions: ", err)
		// TODO: Create Incident
		return
	}

	//TODO: UPDATE CACHE key: /block-bumber/abbr/101010
}

func getTransactions(ctx context.Context, abbr string, lastBlock, lastSavedBlock uint64) ([]*types.Transaction, error) {
	lastBlock = 19440138 // TODO: REMOVE BEFORE DEPLOY
	txs := []*types.Transaction{}

	if lastBlock <= lastSavedBlock {
		return []*types.Transaction{}, nil
	}

	for blockNum := lastSavedBlock + 1; blockNum <= lastBlock; blockNum++ {
		block, err := chainapi.API(abbr).GetBlock(ctx, blockNum)
		if err != nil {
			return []*types.Transaction{}, err
		}

		txs = append(txs, block.Transactions()...)
	}

	return txs, nil
}

func filterTransactions(txs []*types.Transaction) ([]*types.Transaction, error) {
	filteredTxs := []*types.Transaction{}
	for _, tx := range txs {
		var from, to bool
		fromAddress, toAddress, err := chainapi.GetAddresses(tx)
		if err != nil {
			return []*types.Transaction{}, err
		}

		// TODO: GET DATA FROM CACHE
		_, from = mockedAddresses[fromAddress]
		_, to = mockedAddresses[toAddress]

		if to || from {
			filteredTxs = append(filteredTxs, tx)
		}
	}

	return filteredTxs, nil
}

func processTransactions(ctx context.Context, abbr string, txs []*types.Transaction) error {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, tx := range txs {
		value := tx
		eg.Go(func() error { return notify(egCtx, abbr, value) })
	}

	return eg.Wait()
}

func notify(ctx context.Context, abbr string, tx *types.Transaction) error {
	if err := chainapi.API(abbr).ValidateTransaction(ctx, tx); err != nil {
		if errors.Is(err, chainapi.ErrTxFailed) {
			return nil
		}
		return err
	}

	// TODO: SEND MESSAGE TO NOTIFY CONSUMER (topic = transaction-notify) ?

	return nil
}
