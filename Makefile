build:
	protoc -I. --go_out=plugins=micro:. \
	  proto/user/user.proto

	docker build -t user-service .

run:
