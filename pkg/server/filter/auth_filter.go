package filter

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/constants"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
	"time"
)

const (
	tokenHeaderKey = "Authorization"
	tokenCookieKey = "session"
)

var secret = ""

func SetupSecret(sec string) {
	secret = sec
}

func AuthenticateValidate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	uToken, err := extractToken(req.Request)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	token, err := validate(uToken)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	err = injectContext(req, token)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	chain.ProcessFilter(req, resp)
}

func extractToken(r *http.Request) (string, error) {
	jwtHeader := strings.Split(r.Header.Get(tokenHeaderKey), " ")

	if jwtHeader[0] == "Bearer" && len(jwtHeader) == 2 {
		return jwtHeader[1], nil
	}

	jwtCookie, err := r.Cookie(tokenCookieKey)

	if err == nil {
		return jwtCookie.Value, nil
	}

	return "", fmt.Errorf("no token found")
}

func validate(uToken string) (*jwt.Token, error) {
	if len(uToken) == 0 {
		return nil, fmt.Errorf("token length is zero")
	}

	token, err := jwt.Parse(uToken, providerKey)
	if err != nil {
		klog.Error("parse token error is ", err)
		return nil, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		klog.Error("invalid payload")
		return nil, fmt.Errorf("invalid payload")
	}

	if !payload.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("token is expires")
	}

	//TODO: other validate
	return token, nil
}

func providerKey(token *jwt.Token) (interface{}, error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodHMAC:
		return secret, nil
	default:
		return secret, nil
	}
}

func injectContext(req *restful.Request, token *jwt.Token) error {
	payload, _ := token.Claims.(jwt.MapClaims)
	username, ok := payload["sub"].(string)
	if !ok {
		return fmt.Errorf("invalid payload")
	}
	req.SetAttribute(constants.GoCronUsernameHeader, username)
	return nil
}
