package saml

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	"github.com/russellhaering/goxmldsig"
	log "github.com/sirupsen/logrus"
)

var (
	AvailableRoles  = "https://aws.amazon.com/SAML/Attributes/Role"
	SessionName     = "https://aws.amazon.com/SAML/Attributes/RoleSessionName"
	SessionDuration = "https://aws.amazon.com/SAML/Attributes/SessionDuration"
)

type SAMLService struct {
	log *log.Entry
	sp  *saml2.SAMLServiceProvider
}

func NewAWSSAMLService(idpMetadata string, url string) (*SAMLService, error) {
	metadataRaw, err := base64.StdEncoding.DecodeString(idpMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not read IDP metadata: %v", err)
	}

	metadata := types.EntityDescriptor{}
	err = xml.Unmarshal(metadataRaw, &metadata)
	if err != nil {
		return nil, fmt.Errorf("could not parse IDP metadata: %v", err)
	}

	// create memory store and add certs in metadata to it
	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}
	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				return nil, fmt.Errorf("metadata certificate(%d) must not be empty", idx)
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, err
			}
			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, err
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// we sign the AuthnRequest with a random key because Okta doesn't seem to verify these.
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		SignAuthnRequests:           true,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
		AssertionConsumerServiceURL: url,
	}

	return &SAMLService{
		log: log.WithField("service", "saml"),
		sp:  sp,
	}, nil
}

func (s *SAMLService) SAMLRequestURL() (string, error) {
	url, err := s.sp.BuildAuthURL("")
	if err == nil {
		s.log.Debugf("Built SAML request: %s", url)
	}
	return url, err
}

func (s *SAMLService) ParseSAMLResponse(samlResponse string) (*SAMLInfo, error) {
	s.log.Debug("Parsing SAML response")
	assertionInfo, err := s.sp.RetrieveAssertionInfo(samlResponse) // also verifies signature
	if err != nil {
		return nil, err
	}

	if assertionInfo.WarningInfo.InvalidTime {
		return nil, errors.New("invalid time in assertion")
	}

	if assertionInfo.WarningInfo != nil {
		s.log.Warningf("SAML warnings: %v", assertionInfo.WarningInfo)
	}

	sessionInfo, err := CreateSAMLSessionInfo(assertionInfo)
	if err != nil {
		return sessionInfo, err
	}

	sessionInfo.RawSAML = samlResponse
	s.log.Debugf("Got session info from SAML: %v", sessionInfo)
	return sessionInfo, nil
}
