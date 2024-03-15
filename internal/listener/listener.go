package listener

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/MarcoBuarque/chain-alert/chain-alert/constants"
	chainapi "github.com/MarcoBuarque/chain-alert/chain-alert/pkg/chainAPI"
)

var (
	mockedAddresses = map[string]struct{}{
		"0x2B7d796e57B6d0eE7C0852badA76cD219E1ff88f": {}, //tx: 0x8e1341aeb69fbecf8cde9197e87fc4d3d15ceeac4f88315244c1668a182d7c92
		"0x61029169F4F71583A96c9ded93a1d80d586D24B2": {}, // tx: 0x6d3b0e99af0181b1ac01c6bbd752a365cfebffce6702d9dc85fe599cb92ad63c
		"0xd993C3fae4800538FB07945D52C61B55Fc4f9F8c": {}, // tx: 0xc0cb372e73f6c391f7710a147ea2ebb346c51fe18342635420554f3d261616e4
		"0x0505E982ec4D3Fdcb8F99e061dD82399f6842467": {}, // tx: 0x6a79baad4e5ddae96dba43ae238307a2ea1e50f2f7b56e58ccf2dd05c5f84acd
		"0x7BCf7078c4c0AA54E1F808592AE148d4C3A6aF5b": {}, // tx: 0x8e202993dff70cea7beb2fee0e00c020e3456500c098c0abc3c71b0f8cbef09c UNISWAP ERC20
	}

	syncGroup sync.WaitGroup
)

func Run() error {
	timerStr := os.Getenv("timer")

	timer, err := strconv.Atoi(timerStr)
	if err != nil {
		fmt.Println("invalid timer env: ", timerStr)
		timer = 2
	}

	ticker := time.NewTicker(time.Duration(timer) * time.Second)
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	defer chainapi.Close()

	// TODO: POPULATE lastSavedBlock with cache data map[abbr]blockNumber  (review: load in memory?)
	for {
		select {
		case <-done:
			ticker.Stop()
			fmt.Println("Shutting down...")
			return nil

		case t := <-ticker.C:
			fmt.Println("Checking blocks", t)
			for _, abbr := range constants.AVALIABLE_CHAINS_ABBR {
				syncGroup.Add(len(constants.AVALIABLE_CHAINS_ABBR))

				go func(value string) {
					defer syncGroup.Done()
					processBlock(value)
				}(abbr)
			}
			syncGroup.Wait()
		}
	}

}
