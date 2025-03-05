package catalyst

import (
	"context"

	"github.com/ethereum/go-ethereum/params"
)

func (api *ConsensusAPI) GetChainConfig(_ context.Context) (*params.ChainConfig, error) {
	return api.eth.BlockChain().Config(), nil
}
