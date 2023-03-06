FROM golang:1.19.1-alpine
RUN apk update && apk add git
RUN mkdir /go/src/app
WORKDIR /go/src/app
ADD ./app /go/src/app
RUN go mod download
EXPOSE 80
# CMD [ "go", "run", "main.go"]