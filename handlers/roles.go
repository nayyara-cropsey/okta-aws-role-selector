package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"okta-aws-role-selector/saml"
	"strings"
)

const (
	SAMLResponseParam = "SAMLResponse"
	RelayStateParam   = "RelayState"
	DevPrefix         = "Dev/"
)

type genericResponse map[string]string

type RelayState struct {
	AccountID string
	IsDev     bool
}

func NewRelayState(value string) *RelayState {
	relayState := new(RelayState)
	if strings.HasPrefix(value, DevPrefix) {
		relayState.AccountID = strings.ReplaceAll(value, DevPrefix, "")
		relayState.IsDev = true
	} else {
		relayState.AccountID = value
	}
	return relayState
}

func RolesHandler(templatePage string, config *saml.Config) (gin.HandlerFunc, error) {
	samlService, err := saml.NewAWSSAMLService(config.IdpMetadata, config.SpUrl)
	logger := log.WithField("handler", "roles")

	return func(ctx *gin.Context) {
		samlResponse := ctx.PostForm(SAMLResponseParam)
		if samlResponse == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, genericResponse{"message": fmt.Sprintf("no [%s] parameter found", SAMLResponseParam)})
			return
		}

		samlInfo, err := samlService.ParseSAMLResponse(samlResponse)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, genericResponse{"message": fmt.Sprintf("unable to parse assertion: %s", err)})
			return
		}
		config.UpdateMetaData(samlInfo)

		// if relay state is set, narrow down to specific account
		relayState := NewRelayState(ctx.PostForm(RelayStateParam))
		logger.Infof("Relay state: %v", relayState)
		if relayState.AccountID != "" {
			for _, account := range samlInfo.Accounts {
				if relayState.AccountID == account.ID {
					if relayState.IsDev {
						account.Url = config.DevAccountUrls[relayState.AccountID]
					}
					samlInfo.Accounts = []*saml.Account{account}
					break
				}
			}
		}

		ctx.HTML(http.StatusOK, templatePage, gin.H{
			"Data": samlInfo,
		})
	}, err
}
