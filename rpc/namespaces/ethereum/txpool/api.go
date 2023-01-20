package txpool

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/evmos/ethermint/rpc/backend"
	"github.com/evmos/ethermint/rpc/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// PublicAPI offers and API for the transaction pool. It only operates on data that is non-confidential.
// NOTE: For more info about the current status of this endpoints see https://github.com/evmos/ethermint/issues/124
type PublicAPI struct {
	logger  log.Logger
	backend backend.EVMBackend
}

// NewPublicAPI creates a new tx pool service that gives information about the transaction pool.
func NewPublicAPI(logger log.Logger, backend backend.EVMBackend) *PublicAPI {
	return &PublicAPI{
		logger:  logger.With("module", "txpool"),
		backend: backend,
	}
}

// Content returns the transactions contained within the transaction pool
// TODO: replace this
func (api *PublicAPI) Content() (map[string]map[string]map[string]*types.RPCTransaction, error) {
	api.logger.Debug("txpool_content")
	content := map[string]map[string]map[string]*types.RPCTransaction{
		"pending": make(map[string]map[string]*types.RPCTransaction),
		"queued":  make(map[string]map[string]*types.RPCTransaction),
	}
	return content, nil
}

// Inspect returns the content of the transaction pool and flattens it into an
// TODO: replace this
func (api *PublicAPI) Inspect() (map[string]map[string]map[string]string, error) {
	api.logger.Debug("txpool_inspect")
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}

	pending, err := api.backend.PendingTransactions()
	if err != nil {
		return content, nil
	}

	for _, tx := range pending {
		p, err := evmtypes.UnwrapEthereumMsg(tx, common.Hash{})
		if err != nil {
			// not valid ethereum tx
			continue
		}

		content["pending"][p.Hash] = make(map[string]string)
	}

	return content, nil
}

// Status returns the number of pending and queued transaction in the pool.
// TODO: replace this
func (api *PublicAPI) Status() map[string]hexutil.Uint {
	api.logger.Debug("txpool_status")
	pending, err := api.backend.PendingTransactions()
	if err != nil {
		return map[string]hexutil.Uint{
			"pending": hexutil.Uint(0),
			"queued":  hexutil.Uint(0),
			"success": hexutil.Uint(0),
		}
	}

	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(len(pending)),
		"queued":  hexutil.Uint(0),
		"success": hexutil.Uint(1),
	}
}
