package catalyst

import (
	"context"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/params"
)

func (api *ConsensusAPI) GetChainConfig(context.Context) (*params.ChainConfig, error) {
	return api.eth.BlockChain().Config(), nil
}

func (api *ConsensusAPI) GetFullPayload(payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	return api.getPayload(payloadID, true)
}
