package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/opxyc/go-l-grpc-tts/say"

	"google.golang.org/grpc"
)

func main() {
	backend := flag.String("b", "localhost:8080", "addr of the 'say' backend")
	output := flag.String("o", "output.wav", "output file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speak\"\n", os.Args[0])
		os.Exit(0)
	}

	log.Println("Using backend:", *backend)

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{Text: flag.Arg(0)}
	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatal(err)
	}

	log.Println("Saved output to", *output)
}
