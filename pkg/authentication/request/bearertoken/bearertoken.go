package bearertoken

import (
	"errors"
	"github.com/x893675/gocron/pkg/authentication/authenticator"
	"net/http"
	"strings"
)

type Authenticator struct {
	auth authenticator.Token
}

func New(auth authenticator.Token) *Authenticator {
	return &Authenticator{auth}
}

var invalidToken = errors.New("invalid bearer token")

func (a *Authenticator) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
	auth := strings.TrimSpace(req.Header.Get("Authorization"))
	if auth == "" {
		return nil, false, nil
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, false, nil
	}

	token := parts[1]

	// Empty bearer tokens aren't valid
	if len(token) == 0 {
		return nil, false, nil
	}

	resp, ok, err := a.auth.AuthenticateToken(req.Context(), token)
	// if we authenticated successfully, go ahead and remove the bearer token so that no one
	// is ever tempted to use it inside of the API server
	if ok {
		req.Header.Del("Authorization")
	}

	// If the token authenticator didn't error, provide a default error
	if !ok && err == nil {
		err = invalidToken
	}

	return resp, ok, err
}
