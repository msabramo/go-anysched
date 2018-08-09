# build stage
FROM golang:alpine AS build-env
RUN apk add --update git
RUN wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/github.com/msabramo/go-anysched
WORKDIR /go/src/github.com/msabramo/go-anysched 
RUN dep ensure -v
RUN go build -o bin/anysched-cli ./cmd/anysched-cli

# final stage
FROM alpine
COPY --from=build-env /go/src/github.com/msabramo/go-anysched/bin/anysched-cli /bin
WORKDIR /workdir
ENTRYPOINT ["/bin/anysched-cli"]
