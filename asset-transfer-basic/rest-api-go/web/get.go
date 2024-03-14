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

// 	// Get the channel client
// 	ctx := mockChannelProvider("mychannel")

// 	c, err := New(ctx)
// 	if err != nil {
// 		fmt.Println("failed to create client")
// 	}

// 	block, err := c.QueryBlock(1)
// 	if err != nil {
// 		fmt.Printf("failed to query block: %s\n", err)
// 	}

// 	if block != nil {
// 		fmt.Println("Retrieved block #1")
// 	}
// 	channelClient, err := setup.Gateway.GetNetwork(reqBody.ChannelName).BlockEvents(c)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to get channel client: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// // Query blockchain info
// 	// blockchainInfo, err := channelClient.QueryInfo()
// 	// if err != nil {
// 	// 	http.Error(w, fmt.Sprintf("Failed to query blockchain info: %s", err), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// // Fetch blocks from the genesis block to the latest block
// 	// var blocks []struct {
// 	// 	BlockNumber       uint64 `json:"blockNumber"`
// 	// 	PreviousBlockHash string `json:"previousBlockHash"`
// 	// 	DataHash          string `json:"dataHash"`
// 	// 	NumTransactions   int    `json:"numTransactions"`
// 	// }
// 	// for i := uint64(0); i < blockchainInfo.BCI.Height; i++ {
// 	// 	block, err := channelClient.QueryBlock(i)
// 	// 	if err != nil {
// 	// 		http.Error(w, fmt.Sprintf("Failed to query block %d: %s", i, err), http.StatusInternalServerError)
// 	// 		return
// 	// 	}

// 	// 	blocks = append(blocks, struct {
// 	// 		BlockNumber       uint64 `json:"blockNumber"`
// 	// 		PreviousBlockHash string `json:"previousBlockHash"`
// 	// 		DataHash          string `json:"dataHash"`
// 	// 		NumTransactions   int    `json:"numTransactions"`
// 	// 	}{
// 	// 		BlockNumber:       block.Header.Number,
// 	// 		PreviousBlockHash: fmt.Sprintf("%x", block.Header.PreviousHash),
// 	// 		DataHash:          fmt.Sprintf("%x", block.Header.DataHash),
// 	// 		NumTransactions:   len(block.Data.Data),
// 	// 	})
// 	// }

// 	// // Encode blocks to JSON and send as response
// 	// if err := json.NewEncoder(w).Encode(blocks); err != nil {
// 	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	// 	return
// 	// }
// }
