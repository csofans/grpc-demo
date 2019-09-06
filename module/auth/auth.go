package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	pb "pikachu/demo/protobuf/demo"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/status"
)

var tokenByUser map[string]Data

func init() {
	tokenByUser = make(map[string]Data)
}

// Server -
type Server struct {
}

// Data -
type Data struct {
	ID     string
	Token  string
	UserID string
}

// Login -
func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	log.Printf("call login with param : %v  \n", in)
	token := hashToken(in.GetAccount(), in.GetPassword())

	tokenByUser[token] = Data{
		ID:     time.Now().Format(time.RFC3339),
		Token:  token,
		UserID: in.GetAccount(),
	}

	return &pb.LoginResponse{Token: token, Status: &pb.StatusReply{Code: 0, Msg: "Success", Unix: ptypes.TimestampNow()}}, nil
}

// Logout -
func (s *Server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Printf("call logout with param : %v  \n", in)

	resp := &pb.LogoutResponse{}

	if _, ok := tokenByUser[in.GetToken()]; !ok {
		st := status.New(123, "error")
		return resp, st.Err()
	}
	delete(tokenByUser, in.GetToken())
	resp.Status = &pb.StatusReply{
		Code: 0,
		Msg:  "success",
		Unix: ptypes.TimestampNow(),
	}
	return resp, nil
}

// List -
func (s *Server) List(ctx context.Context, in *pb.TokenListRequest) (*pb.TokenListResponse, error) {
	log.Printf("call list with param : %v  \n", in)
	total := int32(len(tokenByUser))

	ans := make([]*pb.TokenContext, 0)

	for _, v := range tokenByUser {
		ans = append(ans, &pb.TokenContext{
			Id:     v.ID,
			Token:  v.Token,
			Userid: v.UserID,
		})
	}

	return &pb.TokenListResponse{
		Total: total,
		Data:  ans,
	}, nil
}

// ListBySteam -
func (s *Server) ListBySteam(in *pb.TokenListRequest, stream pb.Auth_ListBySteamServer) error {
	log.Printf("call listbystream with param : %v  \n", in)

	total := int32(len(tokenByUser))

	for _, v := range tokenByUser {
		stream.Send(&pb.TokenListStreamResponse{
			Total: total,
			Data: &pb.TokenContext{
				Token:  v.Token,
				Id:     v.ID,
				Userid: v.UserID,
			},
		})
	}

	return nil
}

func hashToken(account, password string) string {
	s := md5.New()
	s.Write([]byte(fmt.Sprintf("%v+%v + %v", account, password, time.Now().Unix())))
	return hex.EncodeToString(s.Sum(nil))
}
