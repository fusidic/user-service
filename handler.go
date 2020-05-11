// user-service/handler.go
package main

import (
	"encoding/json"
	"errors"
	"log"

	pb "github.com/fusidic/user-service/proto/user"
	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
	PubSub       broker.Broker
}

func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(ctx, req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	// _, err := srv.repo.GetByEmailAndPassword(req)
	// if err != nil {
	// 	return err
	// }
	// res.Token = "testingabc"
	// return nil

	// using hash to encrypt req
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	log.Println(user)

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	// if err := srv.repo.Create(req); err != nil {
	// 	return err
	// }
	// res.User = req
	// return nil

	// 对密码进行哈希
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPass)
	if err := srv.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	if err := srv.publishEvent(req); err != nil {
		return err
	}
	return nil
}

func (srv *service) publishEvent(user *pb.User) error {
	// Marshal 序列化 JSON 字符串
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 创建创建事件消息
	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}

	// 发布消息到消息代理中
	if err := srv.PubSub.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}

func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// Decode token
	claims, err := srv.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	log.Println(claims)

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
