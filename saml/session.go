package saml

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/russellhaering/gosaml2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	logger = log.WithField("component", "session")
)

func CreateSAMLSessionInfo(assertions *saml2.AssertionInfo) (*SAMLInfo, error) {
	sessionInfo := new(SAMLInfo)
	accountsMap := make(map[string]*Account)

	// available roles
	values := assertions.Values[AvailableRoles].Values
	logger.Infof("Test: %s", values)
	if len(values) > 0 {
		var convertedValues []string
		for _, value := range values {
			convertedValues = append(convertedValues, value.Value)
		}

		availableRoles := convertedValues
		for _, role := range availableRoles {

			roleInfo, err := extractRoleInfo(role)
			if err != nil {
				return nil, err
			}
			logger.Debugf("Processing role: %s",roleInfo)
			if account, ok := accountsMap[roleInfo.AccountID]; ok {
				account.Roles = append(account.Roles, roleInfo)
			} else {
				accountsMap[roleInfo.AccountID] = &Account{
					ID:    roleInfo.AccountID,
					Roles: []*Role{roleInfo},
				}
			}
		}
	} else {
		return nil, errors.New("failed to parse roles info")
	}

	for _, account := range accountsMap {
		sessionInfo.Accounts = append(sessionInfo.Accounts, account)
	}

	// session name
	values = assertions.Values[SessionName].Values
	if len(values) > 0 {
		sessionInfo.SessionName = values[0].Value
	} else {
		return nil, errors.New("failed to parse session name")
	}

	// duration
	values = assertions.Values[SessionDuration].Values
	if len(values) > 0 {
		sessionInfo.SessionDuration, _ = strconv.Atoi(values[0].Value)
	} else {
		return nil, errors.New("failed to parse session duration")
	}

	logger.Debugf("SessionInfo: %s", sessionInfo)
	return sessionInfo, nil
}

func extractRoleInfo(value string) (*Role, error) {
	roleInfo := new(Role)
	tokens := strings.Split(value, ",")
	if len(tokens) != 2 {
		return roleInfo, errors.New("requires a token with 2 parts")
	}

	// Amazon's documentation suggests that the
	// Role ARN should appear first in the comma-delimited
	// set in the Role Attribute that SAML IdP returns.
	//
	// See the section titled "An Attribute element with the Name attribute set
	// to https://aws.amazon.com/SAML/Attributes/Role" on this page:
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_saml_assertions.html
	//
	// In practice, though, Okta SAML integrations with AWS will succeed
	// with either the role or principal ARN first, and these `if` statements
	// allow that behavior in this program.
	if strings.Contains(tokens[0], ":saml-provider/") {
		// if true, Role attribute is formatted like:
		// arn:aws:iam::ACCOUNT:saml-provider/provider,arn:aws:iam::account:role/roleName
		roleInfo.Arn = tokens[1]
	} else if strings.Contains(tokens[1], ":saml-provider/") {
		// if true, Role attribute is formatted like:
		// arn:aws:iam::account:role/roleName,arn:aws:iam::ACCOUNT:saml-provider/provider
		roleInfo.Arn = tokens[0]
	}
	roleArn, _ := arn.Parse(roleInfo.Arn)
	roleInfo.Name = strings.Replace(roleArn.Resource, "role/", "", -1)
	roleInfo.AccountID = roleArn.AccountID

	return roleInfo, nil
}
