// Package cert contains certificate specifications and
// certificate-specific management.
package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

func displayName(name pkix.Name) string {
	var ns []string

	if name.CommonName != "" {
		ns = append(ns, name.CommonName)
	}

	for _, val := range name.Country {
		ns = append(ns, fmt.Sprintf("C=%s", val))
	}

	for _, val := range name.Organization {
		ns = append(ns, fmt.Sprintf("O=%s", val))
	}

	for _, val := range name.OrganizationalUnit {
		ns = append(ns, fmt.Sprintf("OU=%s", val))
	}

	for _, val := range name.Locality {
		ns = append(ns, fmt.Sprintf("L=%s", val))
	}

	for _, val := range name.Province {
		ns = append(ns, fmt.Sprintf("ST=%s", val))
	}

	if len(ns) > 0 {
		return "/" + strings.Join(ns, "/")
	}

	return ""
}

// Compare if hostnames in certificate and spec are equal
func hostnamesMatchesCertificate(hosts []string, cert *x509.Certificate) bool {
	a := make([]string, len(hosts))
	for idx := range hosts {
		// normalize the IPs.
		ip := net.ParseIP(hosts[idx])
		if ip == nil {
			a[idx] = hosts[idx]
		} else {
			a[idx] = ip.String()
		}
	}
	b := make([]string, len(cert.DNSNames), len(cert.DNSNames)+len(cert.IPAddresses))
	copy(b, cert.DNSNames)
	for idx := range cert.IPAddresses {
		b = append(b, cert.IPAddresses[idx].String())
	}

	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)
	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}
	return true
}

func verifyCertChain(ca *x509.Certificate, cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	roots.AddCert(ca)
	_, err := cert.Verify(x509.VerifyOptions{
		Roots: roots,
	})
	return err
}

func encodeKeyToPem(key interface{}) ([]byte, error) {
	switch key.(type) {
	case *ecdsa.PrivateKey:
		data, err := x509.MarshalECPrivateKey(key.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(
			&pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: data,
			},
		), nil
	case *rsa.PrivateKey:
		return pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key.(*rsa.PrivateKey)),
			},
		), nil
	}
	return nil, errors.New("private key is neither ecdsa nor rsa thus cannot be encoded")
}

// encodeCertificateToPEM serialize a certificate into pem format
func encodeCertificateToPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
}

// strictJSONUnmarshal unmarshals a byte source into the given interface while also
// enforcing that there is no unknown fields
func strictJSONUnmarshal(data []byte, object interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	err := dec.Decode(object)
	if err != nil {
		return err
	}
	if dec.More() {
		return errors.New("multiple json objects found, only one is allowed")
	}
	return nil
}
