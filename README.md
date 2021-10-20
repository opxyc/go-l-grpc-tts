# Text to Speech ith gRPC - Go

It's all about writing a very basic gRPC server which accepts text as input and returns corresponsing speech to the client in Go. TTS is done using a software called flite which is available at http://www.speech.cs.cmu.edu/flite/.

The goal is to create a gRPC package named [say](./say) with a single functionality (`rpc Say(Text) returns(Speech)`) of receiving a Text(string) and returning Speech(bytes) and use it to create server and client.

```proto
message Text {
    string text = 1;
}

message Speech {
    bytes audio = 1;
}
```

The server aka backend is a gRPC server listening on port 8080 by default which can be changed via flag `p`. It accepts incoming Text and returns Speech as mentioned above. For conversion of Text to Speech it will utilizing the flite binary, available in a Docker container (Ubuntu) where the server itself will be running, using `os.Exec`.

```sh
flite -t "text to speak" -o outputFile.wav
```
will generate an audio file outputfile.wav. We will save it to a temp file, read it again and returns it in binary to the client.

The code for server is available at [`/backend`](./backend)


The client will make request to the server(which can be specified via flag `-d`; defaults to localhost:8080) with "text to speak" as input, receives the binary response and saves it to a file(flag `-o`; defaults to output.wav).

The code for client is available at [`/client`](./client)

We will then put the backend in a [Kubernetes deployment](./backend/kubernetes.yml) (I'm using micro8ks 🙃).

## How to Run
Creating docker image:
```sh
# inside backend/
make build
```
Adding the image to micro8ks registry:
```sh
make push
```
Creating kubernetes deployment(also NodePort service):
```sh
microk8s.kubectl apply -f kubernetes.yml
```
#### Using the client
Get the services:
```sh
be@mC:~$ kubectl get svc
NAME          TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubernetes    ClusterIP   10.152.183.1     <none>        443/TCP          32h
                               👇                           👇
say-service   NodePort    10.152.183.126   <none>        8080:30080/TCP   3h38m
```
If you are on a linux machine, connect to IP:8080 (IP = 10.152.183.126). If you are using minikube or inside a VM, I think it should be IP:30080 which is the port assigned to NodePort service(specified in yaml).

```sh
# from client/
                          👇 
go run main.go -b 10.152.183.126:8080 -o santa.wav "Santa Exists"
# output:
2021/10/20 19:16:30 Using backend: 10.152.183.126:8080
2021/10/20 19:16:30 Saved output to santa.wav
```

---
### Credits <3
[@francesc](https://twitter.com/francesc) for [justforfunc](https://www.youtube.com/c/justforfunc)