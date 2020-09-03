package restplus

import (
	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

func HandleInternalError(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service internal error: ", err)
	_ = response.WriteServiceError(http.StatusInternalServerError, restful.ServiceError{Code: http.StatusInternalServerError, Message: sanitizer.Replace(err.Error())})
}

func HandleBadRequest(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service bad request error: ", err)
	_ = response.WriteServiceError(http.StatusBadRequest, restful.ServiceError{Code: http.StatusBadRequest, Message: sanitizer.Replace(err.Error())})
}

func HandleExpectedFailed(response *restful.Response, req *restful.Request, err error) {
	klog.Error("handle expected failed: ", err)
	_ = response.WriteServiceError(http.StatusExpectationFailed, restful.ServiceError{Code: http.StatusExpectationFailed, Message: sanitizer.Replace(err.Error())})
}

func HandleNotFound(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service not found error: ", err)
	_ = response.WriteServiceError(http.StatusNotFound, restful.ServiceError{Code: http.StatusNotFound, Message: sanitizer.Replace(err.Error())})
}

func HandleForbidden(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service forbidden error: ", err)
	_ = response.WriteServiceError(http.StatusForbidden, restful.ServiceError{Code: http.StatusForbidden, Message: sanitizer.Replace(err.Error())})
}

func HandleConflict(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service conflict error: ", err)
	_ = response.WriteServiceError(http.StatusConflict, restful.ServiceError{Code: http.StatusConflict, Message: sanitizer.Replace(err.Error())})
}

func HandleUnauthorized(response *restful.Response, req *restful.Request, err error) {
	klog.Error("service unauthorized error: ", err)
	_ = response.WriteServiceError(http.StatusUnauthorized, restful.ServiceError{Code: http.StatusUnauthorized, Message: sanitizer.Replace(err.Error())})
}
