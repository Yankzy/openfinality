package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FabricConnection struct {
	Gateway  *client.Gateway
	Network  *client.Network
	Contract *client.Contract
}

func newGrpcConnection(tlsCertPath string, peerEndpoint string) (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, "")

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	return connection, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func getPrivateKeyPath(keystoreDir string) (string, error) {
	files, err := os.ReadDir(keystoreDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == "" {
			// CA generates keys without extension, often ending in _sk
			return filepath.Join(keystoreDir, file.Name()), nil
		}
	}
	return "", fmt.Errorf("no private key found in %s", keystoreDir)
}

func newIdentity(certPath string, mspID string) (*identity.X509Identity, error) {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return nil, err
	}
	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func newSign() (identity.Sign, error) {
	return nil, nil // not used directly, implemented below
}

func newSigner(keyPath string) (identity.Sign, error) {
	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create sign function: %w", err)
	}
	return sign, nil
}

func InitFabricConnection(mspID, peerEndpoint, tlsCertPath, certPath, keystoreDir string) (*FabricConnection, *grpc.ClientConn, error) {
	clientConn, err := newGrpcConnection(tlsCertPath, peerEndpoint)
	if err != nil {
		return nil, nil, err
	}

	id, err := newIdentity(certPath, mspID)
	if err != nil {
		return nil, nil, err
	}

	keyPath, err := getPrivateKeyPath(keystoreDir)
	if err != nil {
		return nil, nil, err
	}

	sign, err := newSigner(keyPath)
	if err != nil {
		return nil, nil, err
	}

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}

	network := gateway.GetNetwork("afro-settlement")
	contract := network.GetContract("afrorail")

	return &FabricConnection{
		Gateway:  gateway,
		Network:  network,
		Contract: contract,
	}, clientConn, nil
}
