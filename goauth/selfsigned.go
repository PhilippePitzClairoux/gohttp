package goauth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"time"
)

// Generate self-signed certs for debugging purposes
func genRSAKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}

	return privateKey
}

func certificateTemplate() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"gohttp-testing"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("0.0.0.0")},
	}
}

func dnsCertificateTemplate(commonName string, dnsName []string) x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"gohttp-testing"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		//IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:                    dnsName,
		PermittedDNSDomains:         dnsName,
		PermittedDNSDomainsCritical: true,
	}
}

func CreateCertificate() tls.Certificate {
	privateKey := genRSAKey()
	template := certificateTemplate()

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		log.Fatalf("could not parse certificate : %s", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privateKey,
	}
}

func CreateDnsCertificate(commonName string, dnsName []string) tls.Certificate {
	privateKey := genRSAKey()
	template := dnsCertificateTemplate(commonName, dnsName)

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		log.Fatalf("could not parse certificate : %s", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privateKey,
	}
}
