package controller

import (
	"context"
	"grpc-demo/module/auth"
	"log"

	pb "grpc-demo/protobuf/demo"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer -
type AuthServer struct {
}

// Login -
func (s *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	checkMetadata(ctx)

	log.Printf("call login with param : %v  \n", in)

	token, err := auth.Login(in.GetAccount(), in.GetPassword())
	_ = token

	if err != nil {
		return &pb.LoginResponse{}, status.New(codes.NotFound, "User Not Found").Err()
	}

	return &pb.LoginResponse{Token: token, Status: &pb.StatusReply{Code: 0, Msg: "Success", Unix: ptypes.TimestampNow()}},
		status.Error(codes.OK, "success")
}

// Logout -
func (s *AuthServer) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	checkMetadata(ctx)
	log.Printf("call logout with param : %v  \n", in)

	resp := &pb.LogoutResponse{}

	if err := auth.Logout(in.GetToken()); err != nil {
		return resp, status.New(codes.NotFound, err.Error()).Err()
	}

	resp.Status = &pb.StatusReply{
		Code: 0,
		Msg:  "success",
		Unix: ptypes.TimestampNow(),
	}
	return resp, status.Error(codes.OK, "success")

}

// List -
func (s *AuthServer) List(ctx context.Context, in *pb.TokenListRequest) (*pb.TokenListResponse, error) {

	checkMetadata(ctx)
	log.Printf("call list with param : %v  \n", in)

	total, data := auth.List()

	ans := make([]*pb.TokenContext, 0)

	for _, v := range data {
		ans = append(ans, &pb.TokenContext{
			Id:     v.ID,
			Token:  v.Token,
			Userid: v.UserID,
		})
	}

	return &pb.TokenListResponse{
		Total: total,
		Data:  ans,
	}, status.Error(codes.OK, "success")

}

// ListBySteam -
func (s *AuthServer) ListBySteam(in *pb.TokenListRequest, stream pb.Auth_ListBySteamServer) error {
	checkMetadata(stream.Context())
	log.Printf("call listbystream with param : %v  \n", in)

	total, data := auth.List()

	for _, v := range data {
		stream.Send(&pb.TokenListStreamResponse{
			Total: total,
			Data: &pb.TokenContext{
				Token:  v.Token,
				Id:     v.ID,
				Userid: v.UserID,
			},
		})
	}

	return status.Error(codes.OK, "success")
}
