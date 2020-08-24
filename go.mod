module github.com/gonethopper/nethopper

go 1.15

//replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace google.golang.org/grpc => google.golang.org/grpc v1.27.0 // indirect

//replace github.com/golang/protobuf/protoc-gen-go => github.com/golang/protobuf/protoc-gen-go v1.2.0

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/coreos/etcd v3.3.24+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/packr v1.30.1
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/context v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/jmoiron/sqlx v1.2.0
	github.com/klauspost/reedsolomon v1.9.9 // indirect
	github.com/lucas-clemente/quic-go v0.18.0
	github.com/mmcloughlin/avo v0.0.0-20200803215136-443f81d77104 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.3.2 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.etcd.io/etcd v3.3.24+incompatible // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)
