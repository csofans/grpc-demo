package cmd

import (
	"fmt"
	"grpc-demo/controller"
	"grpc-demo/env"
	"log"
	"net"

	pbdemo "grpc-demo/protobuf/demo"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Start gRPC Server on Port : %v", env.Port)
		startgRPC()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&env.Port, "port", "p", "5000", "grpc server port")
}

func startgRPC() error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.Port))
	if err != nil {
		log.Fatalf("start grpc server error : %v", err)
	}
	l := logrus.New().WithField("service", "demo")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(l),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_logrus.StreamServerInterceptor(l),
				grpc_recovery.StreamServerInterceptor(),
			),
		),
	)

	reflection.Register(s)

	// grpc 設定

	// 註冊服務
	pbdemo.RegisterAuthServer(s, &controller.AuthServer{})
	pbdemo.RegisterUserServer(s, &controller.UserServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server : %v", err)
	}

	return nil
}
