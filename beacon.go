package beacon

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// NewWithBaseURIAndAuth returns a new client which will use the provided
// baseURL and obtain tokens from the tokenFactory.
func NewWithBaseURIAndAuth(baseURI string, tokenFactory func() string) BaseClient {

	bc := NewWithBaseURI(baseURI)

	jwt := tokenFactory()

	bc.Authorizer = autorest.NewAPIKeyAuthorizerWithHeaders(map[string]interface{}{
		"Authorization": "Bearer " + jwt,
	})

	bc.ResponseInspector = azure.WithErrorUnlessStatusCode(200, 201)

	return bc
}
