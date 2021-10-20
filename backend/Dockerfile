FROM ubuntu:20.04

RUN apt update 
RUN apt-get install flite -y
ADD app /app
ENTRYPOINT [ "/app" ]