# Text to Speech with gRPC - Go

It's all about writing a very basic gRPC server which accepts text as input and returns corresponsing speech to the client in Go. TTS is done using flite which is available at http://www.speech.cs.cmu.edu/flite.

The goal is to create a package named [say](./say) with a single gRPC functionality (`rpc Say(Text) returns(Speech)`) of receiving a Text(string) and returning Speech(bytes) and use it to create server and client.

```proto
package say;

service TextToSpeech {
    rpc Say(Text) returns(Speech) {}
}

message Text {
    string text = 1;
}

message Speech {
    bytes audio = 1;
}
```

To generate Go code from proto:
```sh
protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative say.proto
```

The server aka backend is a gRPC server listening on port 8080 by default which can be changed via flag `p`. It accepts incoming Text and returns Speech. For conversion of Text to Speech it will be utilizing flite which is written is C. We will work with [`cgo`](https://golang.org/cmd/cgo/) to call functions available in flite from Go.

The code for backend/server is available at [`/backend`](./backend). The server will use the package [`flite`](./flite) for TTS generation.

The client can make request to the server(which can be specified via flag `-d`; defaults to localhost:8080) with "text to speak" as input, receives the binary response and saves it to a file(flag `-o`; defaults to output.wav). The code for client is available at [`/client`](./client)

We will then put the backend in a [Kubernetes deployment](./backend/kubernetes.yml) (I'm using micro8ks 🙃).

## How to Run
#### Preparing the backend
Building docker image labelled `say-backend` and client:
```sh
make build
```

Adding the image to micro8ks registry:
```sh
make push
```

Creating kubernetes deployment(and NodePort service):
```sh
microk8s.kubectl apply -f backend/kubernetes.yml
```
*It's optional to use kubernetes. You can directly run the image `say-backend` as well.*

#### Using the client
Get the services:
```sh
$ kubectl get svc
NAME          TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubernetes    ClusterIP   10.152.183.1     <none>        443/TCP          32h
                               👇                           👇
say-service   NodePort    10.152.183.126   <none>        8080:30080/TCP   3h38m
```
If you are on a linux machine, connect to IP:8080 (IP = 10.152.183.126). If you are using minikube or inside a VM, I think it should be IP:30080 which is the port assigned to NodePort service(specified in yaml).

```sh
# from client/
                          👇 
$ ./client -b 10.152.183.126:8080 -o santa.wav "Santa Exists"
2021/10/20 19:16:30 Using backend: 10.152.183.126:8080
2021/10/20 19:16:30 Saved output to santa.wav
```

---
### Credits <3
[@francesc](https://twitter.com/francesc) for [justforfunc](https://www.youtube.com/c/justforfunc)