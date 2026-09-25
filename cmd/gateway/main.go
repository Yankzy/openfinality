package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var connections = make(map[string]*FabricConnection)

func setupConnections() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	orgsDir := filepath.Join(pwd, "network", "organizations")

	configs := []struct {
		ID           string
		MSP          string
		PeerEndpoint string
		TLSCert      string
		Cert         string
		KeyStore     string
	}{
		{
			ID:           "NigeriaMSP",
			MSP:          "NigeriaMSP",
			PeerEndpoint: "localhost:7051",
			TLSCert:      filepath.Join(orgsDir, "nigeria", "peers", "peer0", "tls", "ca.crt"),
			Cert:         filepath.Join(orgsDir, "nigeria", "users", "Admin@nigeriaorg.openfinality.com", "msp", "signcerts", "cert.pem"),
			KeyStore:     filepath.Join(orgsDir, "nigeria", "users", "Admin@nigeriaorg.openfinality.com", "msp", "keystore"),
		},
		{
			ID:           "GhanaMSP",
			MSP:          "GhanaMSP",
			PeerEndpoint: "localhost:9051",
			TLSCert:      filepath.Join(orgsDir, "ghana", "peers", "peer0", "tls", "ca.crt"),
			Cert:         filepath.Join(orgsDir, "ghana", "users", "Admin@ghanaorg.openfinality.com", "msp", "signcerts", "cert.pem"),
			KeyStore:     filepath.Join(orgsDir, "ghana", "users", "Admin@ghanaorg.openfinality.com", "msp", "keystore"),
		},
		{
			ID:           "MoroccoMSP",
			MSP:          "MoroccoMSP",
			PeerEndpoint: "localhost:11051",
			TLSCert:      filepath.Join(orgsDir, "morocco", "peers", "peer0", "tls", "ca.crt"),
			Cert:         filepath.Join(orgsDir, "morocco", "users", "Admin@moroccoorg.openfinality.com", "msp", "signcerts", "cert.pem"),
			KeyStore:     filepath.Join(orgsDir, "morocco", "users", "Admin@moroccoorg.openfinality.com", "msp", "keystore"),
		},
		{
			ID:           "Adapter",
			MSP:          "NigeriaMSP",
			PeerEndpoint: "localhost:7051",
			TLSCert:      filepath.Join(orgsDir, "nigeria", "peers", "peer0", "tls", "ca.crt"),
			Cert:         filepath.Join(orgsDir, "nigeria", "users", "Adapter@nigeriaorg.openfinality.com", "msp", "signcerts", "cert.pem"),
			KeyStore:     filepath.Join(orgsDir, "nigeria", "users", "Adapter@nigeriaorg.openfinality.com", "msp", "keystore"),
		},
	}

	for _, cfg := range configs {
		conn, _, err := InitFabricConnection(cfg.MSP, cfg.PeerEndpoint, cfg.TLSCert, cfg.Cert, cfg.KeyStore)
		if err != nil {
			log.Printf("Warning: failed to init connection for %s: %v", cfg.ID, err)
			continue
		}
		connections[cfg.ID] = conn
		log.Printf("Initialized connection for %s", cfg.ID)
	}

	if len(connections) == 0 {
		return fmt.Errorf("failed to initialize any fabric connections")
	}

	return nil
}

func getConn(r *http.Request) (*FabricConnection, error) {
	orgID := r.Header.Get("X-AfroRail-Org")
	if orgID == "" {
		orgID = "NigeriaMSP" // default fallback for ease of testing
	}
	conn, ok := connections[orgID]
	if !ok {
		return nil, fmt.Errorf("no connection for org ID %s", orgID)
	}
	return conn, nil
}

func main() {
	if err := setupConnections(); err != nil {
		log.Fatalf("Failed to setup connections: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("POST /v1/settlement-intents", createIntentHandler)
	mux.HandleFunc("GET /v1/settlement-intents/{id}", getIntentHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/accept", acceptIntentHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/reservations", reserveIntentHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/execute", executeIntentHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/proof", proofIntentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting gateway on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func respondError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	errCode := "ERR_UNKNOWN"
	if strings.Contains(err.Error(), "ERR_") {
		parts := strings.SplitN(err.Error(), ":", 2)
		if len(parts) > 1 {
			errCode = strings.TrimSpace(parts[0])
		}
	} else if strings.Contains(err.Error(), "already exists") {
		errCode = "ERR_ALREADY_EXISTS"
	}

	json.NewEncoder(w).Encode(map[string]string{
		"code":    errCode,
		"message": err.Error(),
	})
}

func respondSuccess(w http.ResponseWriter, txID string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"transaction_id": txID,
		"status":         "SUCCESS",
	})
}

func createIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}

	var req struct {
		ID                string `json:"id"`
		CounterpartyMSPID string `json:"counterparty_msp_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	_, commit, err := conn.Contract.SubmitAsync("CreateSettlementIntent", client.WithArguments(req.ID, req.CounterpartyMSPID))
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	status, err := commit.Status()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	if !status.Successful {
		respondError(w, fmt.Errorf("commit failed: %v", status), http.StatusInternalServerError)
		return
	}

	respondSuccess(w, commit.TransactionID(), http.StatusCreated)
}

func getIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}
	
	id := r.PathValue("id")
	result, err := conn.Contract.EvaluateTransaction("QueryIntent", id)
	if err != nil {
		respondError(w, err, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func acceptIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")
	
	_, commit, err := conn.Contract.SubmitAsync("AcceptSettlementIntent", client.WithArguments(id))
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	status, err := commit.Status()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	if !status.Successful {
		respondError(w, fmt.Errorf("commit failed: %v", status), http.StatusInternalServerError)
		return
	}

	respondSuccess(w, commit.TransactionID(), http.StatusOK)
}

func reserveIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")
	
	_, commit, err := conn.Contract.SubmitAsync("ReserveCapacity", client.WithArguments(id))
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	status, err := commit.Status()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	if !status.Successful {
		respondError(w, fmt.Errorf("commit failed: %v", status), http.StatusInternalServerError)
		return
	}

	respondSuccess(w, commit.TransactionID(), http.StatusOK)
}

func executeIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")
	
	_, commit, err := conn.Contract.SubmitAsync("EndorseInstruction", client.WithArguments(id))
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	status, err := commit.Status()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	if !status.Successful {
		respondError(w, fmt.Errorf("commit failed: %v", status), http.StatusInternalServerError)
		return
	}

	respondSuccess(w, commit.TransactionID(), http.StatusOK)
}

func proofIntentHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := getConn(r)
	if err != nil {
		respondError(w, err, http.StatusUnauthorized)
		return
	}
	id := r.PathValue("id")
	
	_, commit, err := conn.Contract.SubmitAsync("CommitSettlementProof", client.WithArguments(id))
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	status, err := commit.Status()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	if !status.Successful {
		respondError(w, fmt.Errorf("commit failed: %v", status), http.StatusInternalServerError)
		return
	}

	respondSuccess(w, commit.TransactionID(), http.StatusOK)
}
