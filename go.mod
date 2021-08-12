module github.com/airkits/nethopper

go 1.16

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

require (
	github.com/airkits/csproto v0.0.0-00010101000000-000000000000
	github.com/airkits/proto v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/gin-gonic/gin v1.7.3
	github.com/go-openapi/spec v0.19.15 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/packr v1.30.1
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.3
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/context v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/klauspost/reedsolomon v1.9.9 // indirect
	github.com/lucas-clemente/quic-go v0.19.2
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mmcloughlin/avo v0.0.0-20201130012700-45c8ae10fd12 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.6.9
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.3.2 // indirect
	github.com/ugorji/go v1.2.6 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.etcd.io/etcd v3.3.25+incompatible
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.0.0-20201202200335-bef1c476418a // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/airkits/csproto => ../csproto/

replace github.com/airkits/proto => ../proto/
