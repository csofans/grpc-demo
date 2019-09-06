package cmd

import (
	"fmt"
	"log"
	"net"
	"pikachu/demo/env"
	"pikachu/demo/module/auth"

	pbdemo "pikachu/demo/protobuf/demo"

	//"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
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

	s := grpc.NewServer()

	pbdemo.RegisterAuthServer(s, &auth.Server{})
	//pbdemo.RegisterUserServer(s, user.Server{})
	//pbdemo.RegisterReportServer(s, report.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server : %v", err)
	}
	return nil
}
