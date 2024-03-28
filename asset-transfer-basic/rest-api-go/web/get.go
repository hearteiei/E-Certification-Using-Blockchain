package web

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func (setup *OrgSetup) GetAllBlocksHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extract channel name from request body
// 	var reqBody struct {
// 		ChannelName string `json:"channelName"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
// 		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Get the network client
// 	network, err := setup.Gateway.GetNetwork(reqBody.ChannelName)
// 	if err != nil {
// 		// Handle the error
// 		http.Error(w, fmt.Sprintf("Failed to get network: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	if network == nil {
// 		// Handle the case where network is nil
// 		http.Error(w, "Network client is nil", http.StatusInternalServerError)
// 		return
// 	}

// 	// Get the ledger client from the network
// 	ledgerClient, err := network.GetLedger()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to get ledger client: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	blockchainInfo, err := ledgerClient.QueryInfo()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to query blockchain info: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	var blocks []struct {
// 		BlockNumber       uint64 `json:"blockNumber"`
// 		PreviousBlockHash string `json:"previousBlockHash"`
// 		DataHash          string `json:"dataHash"`
// 		NumTransactions   int    `json:"numTransactions"`
// 	}

// 	// Fetch blocks from the genesis block to the latest block
// 	for i := uint64(0); i < blockchainInfo.BCI.Height; i++ {
// 		block, err := ledgerClient.QueryBlock(i)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to query block %d: %s", i, err), http.StatusInternalServerError)
// 			return
// 		}

// 		blocks = append(blocks, struct {
// 			BlockNumber       uint64 `json:"blockNumber"`
// 			PreviousBlockHash string `json:"previousBlockHash"`
// 			DataHash          string `json:"dataHash"`
// 			NumTransactions   int    `json:"numTransactions"`
// 		}{
// 			BlockNumber:       block.Header.Number,
// 			PreviousBlockHash: fmt.Sprintf("%x", block.Header.PreviousHash),
// 			DataHash:          fmt.Sprintf("%x", block.Header.DataHash),
// 			NumTransactions:   len(block.Data.Data),
// 		})
// 	}

// 	// Encode blocks to JSON and send as response
// 	if err := json.NewEncoder(w).Encode(blocks); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }
