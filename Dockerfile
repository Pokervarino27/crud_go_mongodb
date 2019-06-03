# Dockerfile References: https://docs.docker.com/engine/reference/builder/

#Start from golang v1.10 base image
FROM golang:1.10

#Administrator
LABEL maintainer="Diego Comihual <dcomihual@imaginex.cl>"

#Set the current working directory insider the container
WORKDIR $GOPATH/src/github.com/pokervarino27/crud_go_mongodb

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

#Download all the dependencies
RUN go get -d -v ./...

#Install the package
RUN go install -v ./...

#This container exposes port 8080 to the outside world
EXPOSE 6767

#Run the executable
CMD ["crud_go_mongodb"]
