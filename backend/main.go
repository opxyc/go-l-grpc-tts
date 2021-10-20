package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"

	pb "github.com/opxyc/go-l-t2s/say"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on port", *port)

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

}

type server struct{}

func (server) Say(ctx context.Context, t *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("failed to create tmp file : %s", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close %s:%s", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", t.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("flight failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file : %s", err)
	}

	return &pb.Speech{Audio: data}, nil
}
