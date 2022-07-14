package factory

import (
	"time"

	sm "github.com/tendermint/tendermint/internal/state"
	"github.com/tendermint/tendermint/internal/test/factory"
	"github.com/tendermint/tendermint/types"
)

func MakeBlocks(n int, state *sm.State, privVal types.PrivValidator) []*types.Block {
	blocks := make([]*types.Block, 0)

	var (
		prevBlock     *types.Block
		prevBlockMeta *types.BlockMeta
	)

	appHeight := byte(0x01)
	for i := 0; i < n; i++ {
		height := int64(i + 1)

		block, parts := makeBlockAndPartSet(*state, prevBlock, prevBlockMeta, privVal, height)
		blocks = append(blocks, block)

		prevBlock = block
		prevBlockMeta = types.NewBlockMeta(block, parts)

		// update state
		state.AppHash = []byte{appHeight}
		appHeight++
		state.LastBlockHeight = height
	}

	return blocks
}

func MakeBlock(state sm.State, height int64, c *types.Commit) *types.Block {
	block, _ := state.MakeBlock(
		height,
		MakeData(factory.MakeTenTxs(state.LastBlockHeight), nil, nil),
		c,
		state.Validators.GetProposer().Address,
	)
	return block
}

func MakeData(txs []types.Tx, evd []types.Evidence, msgs []types.Message) types.Data {
	return types.Data{
		Txs: txs,
		Evidence: types.EvidenceData{
			Evidence: evd,
		},
		Messages: types.Messages{
			MessagesList: msgs,
		},
	}
}

func MakeDataFromTxs(txs []types.Tx) types.Data {
	return types.Data{
		Txs: txs,
	}
}

func makeBlockAndPartSet(state sm.State, lastBlock *types.Block, lastBlockMeta *types.BlockMeta,
	privVal types.PrivValidator, height int64) (*types.Block, *types.PartSet) {

	lastCommit := types.NewCommit(height-1, 0, types.BlockID{}, nil)
	if height > 1 {
		vote, _ := factory.MakeVote(
			privVal,
			lastBlock.Header.ChainID,
			1, lastBlock.Header.Height, 0, 2,
			lastBlockMeta.BlockID,
			time.Now())
		lastCommit = types.NewCommit(vote.Height, vote.Round,
			lastBlockMeta.BlockID, []types.CommitSig{vote.CommitSig()})
	}

	return state.MakeBlock(height, MakeDataFromTxs([]types.Tx{}), lastCommit, state.Validators.GetProposer().Address)
}
