protoc ./pb/chat.proto --go_out=plugins=grpc:.
go build  -o chat *.go
