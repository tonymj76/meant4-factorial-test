pto:
	protoc -I. --go_out=plugins=grpc:. proto/factorial.proto

mockfactorial:
	mockgen -source=proto/factorial.pb.go -destination=mock/factorial_service.go -package=mock

run:
	docker run --name factorial -p 5100:5100 fact

dockerbuild:
	 docker build -t fact .