package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"okta-aws-role-selector/saml"
)

const (
	SelectedRoleParam = "SelectedRoleARN"
)

func ExampleHandler(templatePage string, config *saml.Config) (gin.HandlerFunc, error) {
	samlService, err := saml.NewAWSSAMLService(config.IdpMetadata, config.SpUrl)

	return func(ctx *gin.Context) {
		roleInfo := ctx.PostForm(SelectedRoleParam)
		if roleInfo == "" {
			roleInfo = "UNDEFINED"
		}

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

		ctx.HTML(http.StatusOK, templatePage, gin.H{
			"Data":     samlInfo,
			"SelectedRole": roleInfo,
		})
	}, err
}
