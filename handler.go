// user-service/handler.go
package main

import (
	pb "github.com/fusidic/user-service/proto/user"
	"golang.org/x/net/context"
)

type service struct {
	repo         Repository
	tokenService Authable
}

func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {

}
