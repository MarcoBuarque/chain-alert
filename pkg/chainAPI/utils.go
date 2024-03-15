package chainapi

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrTxFailed = errors.New("transaction failed")
)

func GetAddresses(tx *types.Transaction) (from, to string, err error) {
	if tx.To() != nil {
		to = tx.To().String()
	}

	fromAdd, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return "", "", err
	}

	return fromAdd.String(), to, err
}
