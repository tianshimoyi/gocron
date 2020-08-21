package apiserver

import (
	"context"
	"github.com/emicklei/go-restful"
	coreV1 "github.com/x893675/gocron/internal/apiserver/apis/core/v1"
	systemV1 "github.com/x893675/gocron/internal/apiserver/apis/system/v1"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/config"
	"github.com/x893675/gocron/pkg/server/filter"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"net/http"
)

type APIServer struct {
	// webservice container, where all webservice defines
	container *restful.Container
	// db client
	Db *database.Client
	// config
	Config *config.Config
	// http server
	Server *http.Server
	//
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()
	s.container.DoNotRecover(false)
	s.container.Filter(filter.LogRequestAndResponse)
	s.container.Router(restful.CurlyRouter{})
	s.container.RecoverHandler(func(panicReson interface{}, httpWriter http.ResponseWriter) {
		filter.LogStackOnRecover(panicReson, httpWriter)
	})
	s.InstallAPIs()
	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}
	s.Server.Handler = s.container
	return s.Migration()
}

func (s *APIServer) Run(stopCh <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-stopCh
		_ = s.Server.Shutdown(ctx)
	}()
	klog.V(0).Infof("Start listening on %s", s.Server.Addr)
	var err error
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}
	return err
}

func (s *APIServer) InstallAPIs() {
	urlruntime.Must(coreV1.AddToContainer(s.container, s.Db))
	urlruntime.Must(systemV1.AddToContainer(s.container, s.Db))
}

func (s *APIServer) Migration() error {
	return s.Db.DB().Sync2(new(models.Host))
}
