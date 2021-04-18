# grpc-go-course
Following along with Udemy's hands on gRPC course for Golang by Stephane Maarek.  Includes some changes necessary to work with the downloaded protoc.exe precompiled for Windows.

## go mod

Use go mod to get a go.mod and go.sum files.

```
go mod init
```

Periodically especially after getting or importing libraries, run

```
go mod tidy
```

My VSCode could not locate all modules imported by the server and client code without go mod.

## protoc.exe

The protoc.exe reads in .proto files that contain type, message and service as protobuffer definitions, and generates code for the gRPC APIs.

My OS is Windows so I downloaded the win64.zip that contains the compiled protoc.exe executable.

https://github.com/protocolbuffers/protobuf/releases/tag/v3.15.8

After extracting the zipped download, I moved protoc.exe to a folder that was already on my machine's path (to save me having to add a new location the environmental variable where PATH is defined.) I put it into %GOPATH%/bin.

## generate.sh and the protoc command

Perhaps because I'm using a pre-built protoc.exe, the tutorial's protoc command did not work for me to generate the "Greet" pb.go files.  

The protoc.exe I downloaded was expecting an input path that contained slashes, additional flags, and generated two pb.go files.  One of the files has the structures and functions for gRPC commands, and the other contains the type definitions.

## Latest protobuf libraries

I grabbed the latest protobuf libraries with go get and afterward called "go mod tidy".

```
go get -u github.com/golang/protobuf/proto

go get -u github.com/golang/protobuf/protoc-gen-go

go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

go mod tidy
```


## protoc.exe command that worked for me in Windows

I drew from the examples in the official "Quickstart" tutorial.

https://grpc.io/docs/languages/go/quickstart/

```
protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto
```

The --go_out flag generates in greet.pb.go that contains the go type definitions.

The --go_grpc_out flag generates greet_grpc.pb.go code that defines the gRPC APIs.



