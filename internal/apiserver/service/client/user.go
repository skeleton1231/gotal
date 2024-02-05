package service

import (
	"context"

	pb "github.com/skeleton1231/gotal/internal/apiserver/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcClient struct {
	client pb.UserServiceClient
}

func NewUserGrpcClient(connString string) (*UserGrpcClient, error) {
	conn, err := grpc.Dial(connString, grpc.WithTransportCredentials(insecure.NewCredentials())) // 注意安全性设置
	if err != nil {
		return nil, err
	}
	client := pb.NewUserServiceClient(conn)
	return &UserGrpcClient{client: client}, nil
}

func (c *UserGrpcClient) CreateUser(ctx context.Context, user *pb.User) (*pb.CreateResponse, error) {
	// 这里实现调用gRPC方法
	response, err := c.client.Create(ctx, &pb.CreateRequest{User: user})
	if err != nil {
		return nil, err
	}
	return response, nil
}
