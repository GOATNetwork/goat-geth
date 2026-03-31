package catalyst

import (
	"context"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func (api *ConsensusAPI) GetChainConfig(context.Context) (*params.ChainConfig, error) {
	return api.eth.BlockChain().Config(), nil
}

func (api *ConsensusAPI) GetFullPayload(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	return api.getPayload(payloadID, true, nil, nil)
}

type FullPayload struct {
	Data            engine.ExecutableData `json:"data"`
	VersionedHashes []common.Hash         `json:"versionedHashes,omitempty"`
	BeaconRoot      *common.Hash          `json:"beaconRoot,omitempty"`
	Requests        [][]byte              `json:"requests"` // Dont' use hexutil.Bytes here for better performance
}

func (api *ConsensusAPI) PutFullPayload(ctx context.Context, payload *FullPayload) (engine.PayloadStatusV1, error) {
	return api.newPayload(ctx, payload.Data, payload.VersionedHashes, payload.BeaconRoot, payload.Requests, false)
}
