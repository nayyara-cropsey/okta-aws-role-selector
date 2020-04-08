package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okta-aws-role-selector/saml"
)

const (
	SAMLResponseParam = "SAMLResponse"
	RelayStateParam   = "RelayState"
)

type genericResponse map[string]string

func RolesHandler(templatePage string, config *saml.Config) (gin.HandlerFunc, error) {
	samlService, err := saml.NewAWSSAMLService(config.IdpMetadata, config.SpUrl)

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
		relayState := ctx.PostForm(RelayStateParam)
		if relayState != "" {
			for _, account := range samlInfo.Accounts {
				if relayState == account.ID {
					samlInfo.Accounts = []*saml.Account{account}
					break
				}
			}
		}
		samlInfo.Url = relayState

		ctx.HTML(http.StatusOK, templatePage, gin.H{
			"Data": samlInfo,
		})
	}, err
}
