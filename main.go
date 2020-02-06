package main

import (
	"crypto/tls"
	"fmt"
	"net/rpc"
	"time"

	"github.com/bitmark-inc/bitmarkd/counter"
	"golang.org/x/crypto/sha3"
)

var connectionCountRPC counter.Counter

func main() {
	go TLSListen()

	for {
		time.Sleep(1 * time.Minute)
	}
}

func TLSListen() error {
	// load certificate
	listenAddr := "0.0.0.0" + ":2130"
	tlsConfiguration, _, err := getCertificate("client_rpc", cert, key)
	if nil != err {
		return err
	}
	server := rpc.NewServer()
	l, err := tls.Listen("tcp", listenAddr, tlsConfiguration)
	if err != nil {
		fmt.Printf("rpc server listen error: %v \n", err)
		return err
	}
	_, err = l.Accept()
	if err != nil {
		fmt.Printf("rpc.Server terminated: accept error: %v \n", err)
		return err
	}
	go listenAndServeRPC(l, server, 50)
	return nil
}

func getCertificate(name, certificate, key string) (*tls.Config, [32]byte, error) {
	var fingerprint [32]byte

	keyPair, err := tls.X509KeyPair([]byte(certificate), []byte(key))
	if err != nil {
		fmt.Printf("%s failed to load keypair: %v \n", name, err)
		return nil, fingerprint, err
	}

	tlsConfiguration := &tls.Config{
		Certificates: []tls.Certificate{
			keyPair,
		},
	}

	fingerprint = CertificateFingerprint(keyPair.Certificate[0])

	return tlsConfiguration, fingerprint, nil
}
func CertificateFingerprint(certificate []byte) [32]byte {
	return sha3.Sum256(certificate)
}
