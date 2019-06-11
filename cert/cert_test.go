package cert

import (
	"io/ioutil"
	"strings"
	"testing"
)

var (
	ca1 = &CA{
		Name:        "test_ca",
		Remote:      "test_remote",
		Label:       "test_label",
		Profile:     "test_profile",
		AuthKey:     "1111",
		AuthKeyFile: "test_keyfile",
		File:        nil,
	}
	ca2 = &CA{
		Name:        "test_ca",
		Remote:      "test_remote",
		Label:       "test_label",
		Profile:     "test_profile",
		AuthKey:     "1111",
		AuthKeyFile: "test_keyfile",
		File:        nil,
	}
)

func TestCompareCertificatesNewlines(t *testing.T) {
	var tests = []struct {
		name string
		path string
	}{
		{"inner newlines", "test_files/inner-newlines.pem"},
		{"outer newlines", "test_files/outer-newlines.pem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			googlePem, err := ioutil.ReadFile("test_files/google.pem")
			if err != nil {
				t.Fatal(err)
			}

			newlinePem, err := ioutil.ReadFile(tt.path)
			if err != nil {
				t.Fatal(err)
			}

			ca1.SetPEM(googlePem)
			ca2.SetPEM(newlinePem)
			isSame, err := CompareCertificates(ca1.GetPEM(), ca2.GetPEM())
			if err != nil {
				t.Fatal(err)
			}
			if !isSame {
				t.Fatal("Duplicate certificates do not match")
			}
		})
	}
}

func TestCompareCertificateDifferent(t *testing.T) {
	googlePem, err := ioutil.ReadFile("test_files/google.pem")
	if err != nil {
		t.Fatal(err)
	}

	digicertPem, err := ioutil.ReadFile("test_files/digicert.pem")
	if err != nil {
		t.Fatal(err)
	}

	ca1.SetPEM(googlePem)
	ca2.SetPEM(digicertPem)
	isSame, err := CompareCertificates(ca1.GetPEM(), ca2.GetPEM())
	if err != nil {
		t.Fatal(err)
	}
	if isSame {
		t.Fatal("Different certificates match")
	}
}

func TestCompareCertificateNil(t *testing.T) {
	digicertPem, err := ioutil.ReadFile("test_files/digicert.pem")
	if err != nil {
		t.Fatal(err)
	}

	ca1.SetPEM(digicertPem)
	ca2.SetPEM(nil)
	isSame, err := CompareCertificates(ca1.GetPEM(), ca2.GetPEM())
	if strings.Compare(err.Error(), "Unable to pem decode certificate") != 0 || isSame {
		t.Fatal(err)
	}
}

func TestCompareCertificateEmpty(t *testing.T) {
	googlePem, err := ioutil.ReadFile("test_files/google.pem")
	if err != nil {
		t.Fatal(err)
	}

	ca1.SetPEM(googlePem)
	ca2.SetPEM([]byte{})
	isSame, err := CompareCertificates(ca1.GetPEM(), ca2.GetPEM())
	if strings.Compare(err.Error(), "Unable to pem decode certificate") != 0 || isSame {
		t.Fatal(err)
	}
}