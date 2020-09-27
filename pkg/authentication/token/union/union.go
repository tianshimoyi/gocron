package union

import (
	"context"
	"github.com/x893675/gocron/pkg/authentication/authenticator"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// unionAuthTokenHandler authenticates tokens using a chain of authenticator.Token objects
type unionAuthTokenHandler struct {
	// Handlers is a chain of request authenticators to delegate to
	Handlers []authenticator.Token
	// FailOnError determines whether an error returns short-circuits the chain
	FailOnError bool
}

// New returns a token authenticator that validates credentials using a chain of authenticator.Token objects.
// The entire chain is tried until one succeeds. If all fail, an aggregate error is returned.
func New(authTokenHandlers ...authenticator.Token) authenticator.Token {
	if len(authTokenHandlers) == 1 {
		return authTokenHandlers[0]
	}
	return &unionAuthTokenHandler{Handlers: authTokenHandlers, FailOnError: false}
}

// NewFailOnError returns a token authenticator that validates credentials using a chain of authenticator.Token objects.
// The first error short-circuits the chain.
func NewFailOnError(authTokenHandlers ...authenticator.Token) authenticator.Token {
	if len(authTokenHandlers) == 1 {
		return authTokenHandlers[0]
	}
	return &unionAuthTokenHandler{Handlers: authTokenHandlers, FailOnError: true}
}

// AuthenticateToken authenticates the token using a chain of authenticator.Token objects.
func (authHandler *unionAuthTokenHandler) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	var errlist []error
	for _, currAuthRequestHandler := range authHandler.Handlers {
		info, ok, err := currAuthRequestHandler.AuthenticateToken(ctx, token)
		if err != nil {
			if authHandler.FailOnError {
				return info, ok, err
			}
			errlist = append(errlist, err)
			continue
		}

		if ok {
			return info, ok, err
		}
	}

	return nil, false, utilerrors.NewAggregate(errlist)
}
