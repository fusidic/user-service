// user-service/go.mail
package main

import (
	"fmt"

	pb "github.com/fusidic/user-service/proto/user"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/log"
)

func main() {

	// 创建与数据库的连接
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// 自动将用户数据类型转化为数据库的存储类型
	// 每次服务重启之后，都会检查变动，并将数据迁移
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	// 创建服务
	srv := micro.NewService(

		// 名称必须与protobuf中声明的包名一致
		micro.Name("user"),
		micro.Version("latest"),
	)

	// Init 方法会解析所有命令行参数
	srv.Init()

	pubsub := srv.Server().Options().Broker

	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
