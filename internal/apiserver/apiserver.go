package apiserver

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful"
	coreV1 "github.com/x893675/gocron/internal/apiserver/apis/core/v1"
	systemV1 "github.com/x893675/gocron/internal/apiserver/apis/system/v1"
	"github.com/x893675/gocron/internal/apiserver/models"
	hostImpl "github.com/x893675/gocron/internal/apiserver/models/impl/node"
	"github.com/x893675/gocron/internal/apiserver/service/task"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/config"
	"github.com/x893675/gocron/pkg/server/filter"
	"github.com/x893675/gocron/pkg/server/filter/authn/bearertoken"
	"github.com/x893675/gocron/pkg/server/filter/authnz/skyline"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"time"
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
	// taskService
	TaskService *task.Task
	// jwt secret
	JwtSecret string
	//
	SkylineUrl           string
	SkylineAdminRoleName string
	InitAgentHostPath    string
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()
	s.container.DoNotRecover(false)
	s.container.Filter(filter.LogRequestAndResponse)
	s.container.Router(restful.CurlyRouter{})
	s.container.RecoverHandler(func(panicReson interface{}, httpWriter http.ResponseWriter) {
		filter.LogStackOnRecover(panicReson, httpWriter)
	})
	if s.JwtSecret != "" {
		bearertoken.SetupSecret(s.JwtSecret)
		s.container.Filter(bearertoken.AuthenticateValidate)
	} else if s.SkylineUrl != "" {
		skyline.SetupSecret(s.SkylineUrl, 3*time.Second)
		skyline.SetupAdminRoleName(s.SkylineAdminRoleName)
		s.container.Filter(skyline.AuthnzValidate)
	} else {
		klog.V(2).Infof("jwt secret is null, not init authenticate middleware")
	}
	s.InstallAPIs()
	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}
	s.Server.Handler = s.container
	err := s.Migration()
	if err != nil {
		return err
	}
	if err = s.InitAgentHosts(); err != nil {
		klog.Errorf("init agent hosts error: %v", err)
		return err
	}
	return s.TaskService.Initialize(stopCh)
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
	urlruntime.Must(coreV1.AddToContainer(s.container, s.Db, s.TaskService))
	urlruntime.Must(systemV1.AddToContainer(s.container, s.Db))
}

func (s *APIServer) Migration() error {
	return s.Db.DB().Sync2(new(models.Host), new(models.Task), new(models.TaskLog), new(models.TaskHost))
}

func (s *APIServer) InitAgentHosts() error {
	if s.InitAgentHostPath == "" {
		klog.V(1).Infof("no agent hosts init...")
		return nil
	}
	hostModel := hostImpl.New(s.Db)
	_, total, err := hostModel.List(context.TODO(), models.ListHostParam{
		BaseListParam: models.BaseListParam{
			Offset: 0,
			Limit:  10,
		},
	})
	if err != nil {
		return err
	}
	if total > 0 {
		klog.V(1).Infof("hosts has already exist, not init again...")
		return nil
	}
	data, err := readHostsData(s.InitAgentHostPath)
	if err != nil {
		return err
	}
	if err = hostModel.BatchInsert(context.TODO(), data); err != nil {
		return err
	}
	return nil
}

func readHostsData(path string) ([]*models.Host, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data []*models.Host
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
