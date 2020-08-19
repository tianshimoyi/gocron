package restplus

import (
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/pkg/utils/validate"
)

func ParseBody(req *restful.Request, entityPointer interface{}) error {
	err := req.ReadEntity(entityPointer)
	if err != nil {
		return err
	}
	return validate.Validate(entityPointer)
}