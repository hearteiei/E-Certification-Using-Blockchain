package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/rs/cors"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {

	http.HandleFunc("/query", setups.Query)
	http.HandleFunc("/invoke", setups.Invoke)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/checkotp", checkOTPValidity)
	http.HandleFunc("/getall", getAllUsers)
	http.HandleFunc("/generate-certificates", GenerateCertificates)
	corsHandler := cors.Default().Handler(http.DefaultServeMux)
	fmt.Println("Listening (http://localhost:8000/)...")
	if err := http.ListenAndServe(":8000", corsHandler); err != nil {
		fmt.Println(err)
	}
}
