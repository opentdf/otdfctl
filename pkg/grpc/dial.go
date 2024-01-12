package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Conn *grpc.ClientConn
var Context context.Context

func Connect(host string) error {
	var err error
	Conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	Context = context.Background()
	return err
}
