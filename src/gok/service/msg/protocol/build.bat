::protoc --gogofast_out=plugins=grpc:. *.proto

set GOGOPATH=G:\Library\src;F:\alienslib\src
protoc --proto_path=./ --proto_path=%GOPATH% --proto_path=%GOGOPATH% --gogofast_out=plugins=grpc:. *.proto