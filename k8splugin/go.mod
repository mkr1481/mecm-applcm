module k8splugin

go 1.14

require (
	github.com/astaxie/beego v1.12.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/protobuf v1.4.2
	github.com/lib/pq v1.3.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.4.0
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
	helm.sh/helm/v3 v3.3.0
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/metrics v0.18.4
	rsc.io/letsencrypt v0.0.3 // indirect
)
