module github.com/x893675/gocron

go 1.14

require (
	github.com/coreos/go-semver v0.3.0
	github.com/emicklei/go-restful v2.9.6+incompatible
	github.com/emicklei/go-restful-openapi v1.4.1
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.3.3
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/moby/term v0.0.0-20200611042045-63b9a826fb74
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	google.golang.org/grpc v1.27.1
	k8s.io/apimachinery v0.18.8
	k8s.io/component-base v0.18.8
	k8s.io/klog/v2 v2.3.0
	xorm.io/core v0.7.2-0.20190928055935-90aeac8d08eb
)

replace github.com/emicklei/go-restful => github.com/emicklei/go-restful v2.12.0+incompatible
