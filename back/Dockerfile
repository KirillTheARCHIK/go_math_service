FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
  --mount=type=cache,mode=0755,target=/root/.cache/go-build \
  go build -v main.go

CMD [ "go", "run", "." ]