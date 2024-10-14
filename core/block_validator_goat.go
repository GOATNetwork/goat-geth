package core

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

func (v *BlockValidator) validateGoatBlock(block *types.Block) error {
	if v.config.Goat == nil {
		return nil
	}

	extra := block.Header().Extra
	if len(extra) != params.GoatHeaderExtraLengthV0 {
		return fmt.Errorf("no goat tx root found (block %x)", block.Number())
	}

	txLen, txRoot := int(extra[0]), common.BytesToHash(extra[1:])
	if l := block.Transactions().Len(); l < txLen {
		return fmt.Errorf("txs length(%d) is less than goat tx length %d", l, txLen)
	}
	if hash := types.DeriveSha(block.Transactions()[:txLen], trie.NewStackTrie(nil)); hash != txRoot {
		return fmt.Errorf("goat tx root hash mismatch (header value %x, calculated %x)", txRoot, hash)
	}
	if len(block.Withdrawals()) > 0 {
		return errors.New("withdrawals not allowed for goat-geth")
	}

	for i, tx := range block.Transactions() {
		if i < txLen {
			if !tx.IsGoatTx() {
				return fmt.Errorf("transaction %d should be goat tx", i)
			}
			if tx.To() == nil {
				return fmt.Errorf("goat tx %d should have a to address", i)
			}
		} else {
			if tx.IsGoatTx() {
				return fmt.Errorf("transaction %d should not be goat tx", i)
			}
			if tx.Type() == types.BlobTxType {
				return fmt.Errorf("blob transaction %d is not allowed", i)
			}
		}
	}

	return nil
}
