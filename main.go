package main

import (
	"database/sql"
	"easybank/api"
	db "easybank/db/sqlc"
	"easybank/gapi"
	"easybank/pb"
	"easybank/util"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	util.LoadConfig(".")
}

func main() {
	conn, err := sql.Open(util.C.DBDriver, util.C.DBsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)

	go runHttpServer(util.C, store)
	runGrpcServer(util.C, store)
}

func runHttpServer(config *util.Config, store db.Store) {
	server, err := api.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create http server:", err)
	}
	log.Fatalln(server.Start(config.HttpAddress))
}

func runGrpcServer(config *util.Config, store db.Store) {
	server, err := gapi.NewServer(*config, store)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEasybankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcAddress)
	if err != nil {
		log.Fatalln("cannot create listener", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}
}
