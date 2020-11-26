# nethopper
micro-module framework


## mac protobuf 安装
```
brew install protobuf

go get -u -v github.com/golang/protobuf/proto
go get -u -v github.com/golang/protobuf/protoc-gen-go
```



## etcd 和 grpc版本兼容问题

```
go.mod设置grpc版本
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0 

指定protoc版本
 go get github.com/golang/protobuf/protoc-gen-go@v1.3.2 

```