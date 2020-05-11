module gitlab.com/wiserskills/v3/services/userquery

go 1.12

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/arangodb/go-driver v0.0.0-20191207080609-379140cb2abb
	github.com/asdine/storm v2.1.2+incompatible
	github.com/asdine/storm/v3 v3.0.0
	github.com/beevik/etree v1.1.0 // indirect
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgraph-io/dgo v1.0.0
	github.com/dgraph-io/ristretto v0.0.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dimfeld/httptreemux v5.0.1+incompatible // indirect
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/gavv/httpexpect v2.0.0+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jinzhu/gorm v1.9.10
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.7
	github.com/mdempsky/gocode v0.0.0-20190203001940-7fb65232883f // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/nats-io/go-nats v1.7.2
	github.com/nats-io/go-nats-streaming v0.4.4
	github.com/nats-io/nats-streaming-server v0.15.1
	github.com/nats-io/nats.go v1.9.1
	github.com/nats-io/stan.go v0.5.2
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.2.1
	github.com/russellhaering/gosaml2 v0.3.1
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7
	github.com/stretchr/testify v1.3.0
	github.com/unidoc/unioffice v1.2.1 // indirect
	github.com/valyala/fasthttp v1.5.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.10.0
	goa.design/goa v2.0.0+incompatible
	goa.design/goa/v3 v3.0.9
	goa.design/plugins/v3 v3.0.9
	//goa.design/plugins/v3 v3.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	google.golang.org/grpc v1.24.0
)

//replace goa.design/goa/v3 v3.0.9 => ../../../../../goa.design/goa

//replace //goa.design/plugins/v3 v3.0.0 => ../../../../../goa.design/plugins
