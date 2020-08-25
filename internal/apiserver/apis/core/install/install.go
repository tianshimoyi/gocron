package install

import (
	"github.com/emicklei/go-restful"
	v1 "github.com/x893675/gocron/internal/apiserver/apis/core/v1"
	"github.com/x893675/gocron/pkg/server/runtime"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	Install(runtime.Container)
}

func Install(c *restful.Container) {
	urlruntime.Must(v1.AddToContainer(c, nil, nil))
}
