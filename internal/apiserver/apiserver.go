package apiserver

import (
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/config"
)

type APIServer struct {
	// webservice container, where all webservice defines
	container *restful.Container
	// db client
	Db *database.Client
	// config
	Config *config.Config
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	return nil
}

func (s *APIServer) Run(stopCh <-chan struct{}) error {
	return nil
}
