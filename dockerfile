#  syntax=docker/dockerfile:1

# Fetch the base image
FROM golang:1.16-alpine 

# create a folder
WORKDIR /webapp

# copy files 
COPY . /webapp
COPY . ./

WORKDIR /webapp/review

# install go dependencies
RUN go mod download

# build the project
RUN go build -o /web-app ./cmd/web

# expose the project to the public
EXPOSE 4000

# run the project
CMD [ "/web-app" ]
