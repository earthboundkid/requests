package requests

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/vkuznet/x509proxy"
)

// X509HttpClient represents Http client with full support for X509
// (proxy) certificates.

// client X509 certificates
func tlsCerts(key, cert string) ([]tls.Certificate, error) {
	uproxy := os.Getenv("X509_USER_PROXY")
	uckey := os.Getenv("X509_USER_KEY")
	ucert := os.Getenv("X509_USER_CERT")
	if key != "" {
		uckey = key
	}
	if cert != "" {
		ucert = cert
	}

	// check if /tmp/x509up_u$UID exists, if so setup X509_USER_PROXY env
	u, err := user.Current()
	if err == nil {
		fname := fmt.Sprintf("/tmp/x509up_u%s", u.Uid)
		if _, err := os.Stat(fname); err == nil {
			uproxy = fname
		}
	}

	if uproxy == "" && uckey == "" { // user doesn't have neither proxy or user certs
		return nil, nil
	}
	if uproxy != "" {
		// use local implementation of LoadX409KeyPair instead of tls one
		x509cert, err := x509proxy.LoadX509Proxy(uproxy)
		if err != nil {
			return nil, fmt.Errorf("failed to parse X509 proxy: %v", err)
		}
		certs := []tls.Certificate{x509cert}
		return certs, nil
	}
	x509cert, err := tls.LoadX509KeyPair(ucert, uckey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user X509 certificate: %v", err)
	}
	certs := []tls.Certificate{x509cert}
	return certs, nil
}

// X509HttpClient is HTTP client for urlfetch server
func X509HttpClient(key, cert, caPath string, skipVerify bool, verbose bool) *http.Client {
	var certs []tls.Certificate
	var err error
	// get X509 certs
	certs, err = tlsCerts(key, cert)
	if err != nil {
		log.Fatal("ERROR ", err.Error())
	}
	if verbose {
		fmt.Printf("read %d certificates\n", len(certs))
	}
	if len(certs) == 0 {
		return &http.Client{}
	}
	var tr *http.Transport
	if caPath != "" {
		rootCAs := x509.NewCertPool()
		files, err := ioutil.ReadDir(caPath)
		if err != nil {
			log.Fatalf("Unable to list files in '%s', error: %v\n", caPath, err)
		}
		for _, finfo := range files {
			fname := fmt.Sprintf("%s/%s", caPath, finfo.Name())
			caCert, err := os.ReadFile(filepath.Clean(fname))
			if err != nil {
				log.Printf("Unable to read %s\n", fname)
			}
			if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
				if strings.HasSuffix(fname, "pem") {
					log.Printf("invalid PEM format while importing trust-chain: %q", fname)
				}
			}
			if verbose {
				fmt.Printf("read %s\n", fname)
			}
		}
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       certs,
				RootCAs:            rootCAs,
				InsecureSkipVerify: skipVerify},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       certs,
				InsecureSkipVerify: skipVerify},
		}
	}
	if verbose {
		fmt.Printf("HTTP transport %+v\n", tr)
	}
	return &http.Client{Transport: tr}
}
