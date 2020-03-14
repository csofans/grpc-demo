package controller

import (
	"context"
	pb "grpc-demo/protobuf/demo"
	"log"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var usermap map[string]*pb.UserInfo

// UserServer -
type UserServer struct {
}

func init() {
	usermap = make(map[string]*pb.UserInfo)
}

func checkMetadata(ctx context.Context) {

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println("get metadata ok")
	}
	spew.Dump(md)

}

// Create -
func (s *UserServer) Create(ctx context.Context, in *pb.UserInfo) (*pb.StatusReply, error) {
	checkMetadata(ctx)

	if _, ok := usermap[in.GetAccount()]; ok {
		return &pb.StatusReply{}, status.New(codes.AlreadyExists, "user is exists").Err()
	}

	spew.Dump(in.GetRole())

	usermap[in.GetAccount()] = in

	return &pb.StatusReply{}, status.Error(codes.OK, "success")

}

// Delete -
func (s *UserServer) Delete(ctx context.Context, in *pb.UserInfo) (*pb.StatusReply, error) {
	resp := &pb.StatusReply{}
	checkMetadata(ctx)

	if in.GetAccount() == "" {
		return resp, status.New(codes.InvalidArgument, "no account value in the request").Err()
	}

	if _, ok := usermap[in.GetAccount()]; !ok {
		return resp, status.New(codes.NotFound, "no found the user").Err()
	}

	delete(usermap, in.GetAccount())

	return resp, status.Error(codes.OK, "success")

}

// Get -
func (s *UserServer) Get(ctx context.Context, in *pb.UserInfo) (*pb.UserInfoResponse, error) {
	checkMetadata(ctx)
	resp := &pb.UserInfoResponse{}

	if in.GetAccount() == "" {
		return resp, status.New(codes.InvalidArgument, "no account value in the request").Err()
	}

	temp, ok := usermap[in.GetAccount()]
	if !ok {
		return resp, status.New(codes.NotFound, "no found the user").Err()
	}
	resp.User = temp

	return resp, status.Error(codes.OK, "success")
}

// List -
func (s *UserServer) List(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	checkMetadata(ctx)
	resp := &pb.UserListResponse{}

	resp.Total = int32(len(usermap))
	resp.User = func() []*pb.UserInfo {
		x := make([]*pb.UserInfo, len(usermap))
		i := 0
		for _, v := range usermap {
			x[i] = v
		}
		return x
	}()

	return resp, status.Error(codes.OK, "success")
}

// ListSteam -
func (s *UserServer) ListSteam(in *pb.UserListRequest, stream pb.User_ListSteamServer) error {
	checkMetadata(stream.Context())

	total := int32(len(usermap))
	for _, v := range usermap {
		stream.Send(&pb.UserListSteamResponse{
			Total: total,
			User:  v,
		})
	}
	return status.Error(codes.OK, "success")
}
