// user-service/handler.go
package main

import (
	"errors"
	"log"

	pb "github.com/fusidic/user-service/proto/user"
	"github.com/micro/go-micro"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
	Publisher    micro.Publisher
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

	// 新的 publisher 代码更简洁
	if err := srv.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	if err := srv.Publisher.Publish(ctx, req); err != nil {
		return err
	}
	return nil
}

// micro.Publisher 中已经实现了消息的发布，无需此函数了
// publishEvent 中需要手动定义消息 broker.Message
// 并将其进行序列化，造成了额外的开销

// func (srv *service) publishEvent(user *pb.User) error {
// 	// Marshal 序列化 JSON 字符串
// 	body, err := json.Marshal(user)
// 	if err != nil {
// 		return err
// 	}

// 	// 创建创建事件消息
// 	msg := &broker.Message{
// 		Header: map[string]string{
// 			"id": user.Id,
// 		},
// 		Body: body,
// 	}

// 	// 发布消息到消息代理中
// 	if err := srv.PubSub.Publish(topic, msg); err != nil {
// 		log.Printf("[pub] failed: %v", err)
// 	}

// 	return nil
// }

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
