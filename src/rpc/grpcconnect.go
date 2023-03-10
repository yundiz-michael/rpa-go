package rpc

import (
	"context"
	"fmt"
	"github.com/jolestar/go-commons-pool/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// GrpcConnect 连接实例
type GrpcConnect struct {
	IP       string
	Instance *grpc.ClientConn
}

func (connect *GrpcConnect) IsActive() bool {
	state := connect.Instance.GetState()
	return state < 3
}

/*针对不同的ip地址，记录连接池*/
var grpcConnectPoolMap = make(map[string]*pool.ObjectPool)

/*获取ip的grpc连接*/
func getGrpcConnect(ip string) *GrpcConnect {
	ctx := context.Background()
	if _pool, ok := grpcConnectPoolMap[ip]; ok {
		obj, _ := _pool.BorrowObject(ctx)
		return obj.(*GrpcConnect)
	} else {
		factory := pool.NewPooledObjectFactory(
			func(context.Context) (interface{}, error) {
				uri := fmt.Sprintf("%s:3388", ip)
				conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Fatal(err)
				}
				result := &GrpcConnect{
					IP:       ip,
					Instance: conn,
				}
				return result, nil
			}, nil, func(ctx context.Context, object *pool.PooledObject) bool {
				conn := object.Object.(*GrpcConnect)
				return conn.IsActive()
			}, nil, nil)
		ctx := context.Background()
		pool := pool.NewObjectPoolWithDefaultConfig(ctx, factory)
		pool.Config.MaxTotal = 100
		pool.Config.TestOnReturn = true
		pool.Config.TestOnCreate = true
		pool.Config.TestOnBorrow = true
		pool.Config.TestWhileIdle = true
		grpcConnectPoolMap[ip] = pool
		obj, _ := pool.BorrowObject(ctx)
		return obj.(*GrpcConnect)
	}
}

func retGrpcConnect(connect *GrpcConnect) {
	ctx := context.Background()
	if pool, ok := grpcConnectPoolMap[connect.IP]; ok {
		pool.ReturnObject(ctx, connect)
	}
}
