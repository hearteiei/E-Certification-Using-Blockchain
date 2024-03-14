package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() err: %s", err), http.StatusBadRequest)
		return
	}

	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txnProposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating txn proposal: %s", err), http.StatusInternalServerError)
		return
	}
	txnEndorsed, err := txnProposal.Endorse()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error endorsing txn: %s", err), http.StatusInternalServerError)
		return
	}
	txnCommitted, err := txnEndorsed.Submit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error submitting transaction: %s", err), http.StatusInternalServerError)
		return
	}
	// Construct response JSON
	response := map[string]string{
		"transaction_id": txnCommitted.TransactionID(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
