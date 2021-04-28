FROM golang:1.13-alpine as build

ARG _BUILD_FILE

WORKDIR /src

RUN apk update && apk add --no-cache git tzdata && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" >  /etc/timezone

ADD . /src/

RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -ldflags="-w -s" -i -o $_BUILD_FILE cmd/rest/main.go

# user group
RUN echo 'nobody:x:65534:' > /src/group.nobody && \
    echo 'nobody:x:65534:65534::/:' > /src/passwd.nobody

FROM gcr.io/distroless/static
#FROM alpine

ARG _BUILD_FILE

WORKDIR /go/bin

ENV PORT=8080
EXPOSE 8080

# Copy group
COPY --from=build /src/group.nobody /etc/group
COPY --from=build /src/passwd.nobody /etc/passwd
USER nobody:nobody

# Add main program
COPY --from=build /src/$_BUILD_FILE app

ENTRYPOINT ["/go/bin/app"]