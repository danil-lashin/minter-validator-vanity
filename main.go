package main

import (
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/privval"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

func main() {
	pattern := os.Args[1]
	re := regexp.MustCompile(pattern)

	fmt.Printf("Starting search with pattern %s\n", pattern)
	i := int64(0)
	for j := 0; j < 8; j++ {
		go func() {
			for {
				atomic.AddInt64(&i, 1)
				privKey := ed25519.GenPrivKey()

				keyStr := fmt.Sprintf("%x", privKey.PubKey().Bytes()[5:])
				if re.MatchString(keyStr) {
					pv := privval.FilePVKey{
						Address: privKey.PubKey().Address(),
						PubKey:  privKey.PubKey(),
						PrivKey: privKey,
					}
					cdc := amino.NewCodec()
					cryptoAmino.RegisterAmino(cdc)
					jsonBytes, err := cdc.MarshalJSONIndent(pv, "", "  ")
					if err != nil {
						panic(err)
					}
					fmt.Printf("Found: Mp%s\n %s\n", keyStr, jsonBytes)
					os.Exit(0)
				}
			}
		}()
	}

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			<-ticker.C
			fmt.Printf("Checked %d variants\n", i)
		}
	}()

	select {} // run forever
}
