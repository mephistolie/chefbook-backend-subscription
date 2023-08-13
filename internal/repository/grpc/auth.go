package grpc

import (
	api "github.com/mephistolie/chefbook-backend-auth/api/proto/implementation/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Auth struct {
	api.AuthServiceClient
	Conn *grpc.ClientConn
}

func NewAuth(addr string) (*Auth, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, opts)
	if err != nil {
		return nil, err
	}
	return &Auth{
		AuthServiceClient: api.NewAuthServiceClient(conn),
		Conn:              conn,
	}, nil
}
