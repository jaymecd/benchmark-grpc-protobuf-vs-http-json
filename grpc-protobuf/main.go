package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/mail"

	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/grpc-protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lis, _ := net.Listen("tcp", ":60000")

	srv := grpc.NewServer()
	proto.RegisterAPIServer(srv, &Server{})
	log.Println(srv.Serve(lis))
}

type Server struct{}

func (s *Server) CreateUser(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	validationErr := validate(in)
	if validationErr != nil {
		return &proto.Response{
			Code:    500,
			Message: validationErr.Error(),
		}, validationErr
	}

	return &proto.Response{
		Code: 200,
		Id:   "1000000",
	}, nil
}

func validate(in *proto.Request) error {
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		return err
	}

	if len(in.Name) < 4 {
		return errors.New("Name is too short")
	}

	if len(in.Password) < 4 {
		return errors.New("Password is too weak")
	}

	return nil
}
