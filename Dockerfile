# install flite and build the app
FROM golang as builder
RUN apt update && apt install flite-dev -y
RUN mkdir /go/src/app
COPY ./ /go/src/app
WORKDIR /go/src/app/backend
RUN go build -o app

# copy the build to another image 
# (to reduce size of final image)
FROM debian:jessie-slim
RUN apt update && apt install flite-dev -y
COPY --from=builder /go/src/app/backend/app /app
ENTRYPOINT [ "/app" ]