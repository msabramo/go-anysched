# build stage
FROM golang:alpine AS build-env
RUN apk add --update git
RUN wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/git.corp.adobe.com/abramowi/hyperion
WORKDIR /go/src/git.corp.adobe.com/abramowi/hyperion 
RUN dep ensure -v
RUN go build -o bin/hyperion-cli ./cmd/hyperion-cli

# final stage
FROM alpine
COPY --from=build-env /go/src/git.corp.adobe.com/abramowi/hyperion/bin/hyperion-cli /bin
WORKDIR /workdir
ENTRYPOINT ["/bin/hyperion-cli"]
