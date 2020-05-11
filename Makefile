build:
	protoc -I. --go_out=plugins=micro:. \
	  proto/user/user.proto

	docker build -t user-service .

run:
	docker run --net="host" \
		-e MICRO_REGISTRY=mdns \
		-e MICRO_ADDRESS=:50053 \
		-e DB_HOST=localhost \
		-e DB_USER=postgres \
		-e DB_NAME=postgres \
		-e DB_PASSWORD=postgres \
		-e MICRO_BROKER=nats \
		-e MICRO_BROKER_ADDRESS=0.0.0.0:4222 \
		user-service