package skyline

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/constants"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
	"time"
)

const (
	tokenCookieKey = "session"
)

var (
	skylineUrl             = ""
	requestTimeout         = time.Duration(0)
	createHttpRequestError = errors.New("create http request error")
	unauthorizedError      = errors.New("user authorized")
	normalPolicy           = map[string]struct{}{
		"normal:/api/system/v1:GET":  {},
		"normal:/api/system/v1:HEAD": {},
		"normal:/api/core/v1:GET":    {},
		"normal:/api/core/v1:DELETE": {},
		"normal:/api/core/v1:PATCH":  {},
		"normal:/api/core/v1:POST":   {},
	}
	roleAdmin = ""
)

func SetupSecret(url string, timeout time.Duration) {
	skylineUrl = url
	requestTimeout = timeout
}

func SetupAdminRoleName(name string) {
	roleAdmin = name
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain Domain `json:"domain"`
}

type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain Domain `json:"domain"`
}

type Profile struct {
	KeystoneToken string  `json:"keystone_token"`
	Region        string  `json:"region"`
	Project       Project `json:"project"`
	User          User    `json:"user"`
	Roles         []Role  `json:"roles"`
	Exp           int64   `json:"exp"`
	UUID          string  `json:"uuid"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthnzValidate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	uToken, err := extractToken(req.Request)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	profile, err := validate(uToken)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	if !authz(req, profile) {
		restplus.HandleForbidden(resp, req, errors.New("forbidden"))
		return
	}
	err = injectContext(req, profile)
	if err != nil {
		restplus.HandleUnauthorized(resp, req, err)
		return
	}
	chain.ProcessFilter(req, resp)
}

func extractToken(r *http.Request) (string, error) {

	jwtCookie, err := r.Cookie(tokenCookieKey)

	if err == nil {
		return jwtCookie.Value, nil
	}

	return "", fmt.Errorf("no token found")
}

func validate(token string) (*Profile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/profile", skylineUrl), nil)
	if err != nil {
		return nil, createHttpRequestError
	}
	req.Header.Add("Cookie", fmt.Sprintf("%s=%s", tokenCookieKey, token))
	client := &http.Client{}
	if requestTimeout > 0 {
		client.Timeout = requestTimeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusUnprocessableEntity:
		return nil, unauthorizedError
	case http.StatusOK:
		profile := Profile{}
		err = json.NewDecoder(resp.Body).Decode(&profile)
		if err != nil {
			return nil, err
		}
		return &profile, nil
	default:
		return nil, errors.New("unknown skyline status code")
	}
}

func isAdmin(profile *Profile) bool {
	for _, role := range profile.Roles {
		if role.Name == roleAdmin {
			return true
		}
	}
	return false
}

func authz(req *restful.Request, profile *Profile) bool {
	ok := isAdmin(profile)
	if ok {
		req.SetAttribute(constants.GoCronUserRole, constants.UserAdmin)
		return true
	}
	currentParts := splitPath(req.Request.URL.Path)
	if len(currentParts) < 4 {
		klog.Warning("unexpected error, request path should not less than 4 parts")
		req.SetAttribute(constants.GoCronUserRole, constants.UserNormal)
		return false
	}
	policy := fmt.Sprintf("normal:/%s/%s/%s:%s", currentParts[0], currentParts[1], currentParts[2], req.Request.Method)
	if _, ok := normalPolicy[policy]; ok {
		req.SetAttribute(constants.GoCronUserRole, constants.UserAdmin)
		return true
	} else {
		req.SetAttribute(constants.GoCronUserRole, constants.UserNormal)
		return false
	}
}

func injectContext(req *restful.Request, profile *Profile) error {
	req.SetAttribute(constants.GoCronUsernameHeader, profile.User.Name)
	return nil
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
