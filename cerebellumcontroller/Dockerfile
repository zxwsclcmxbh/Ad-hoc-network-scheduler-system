FROM seekwe/golang-gcc:latest as build

RUN mkdir /code

WORKDIR /code

ADD . /code

RUN GOPROXY=https://goproxy.cn,direct go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o main main.go

FROM alpine as prod


COPY --from=build /code/main /

EXPOSE 3000

CMD ["/main"]
