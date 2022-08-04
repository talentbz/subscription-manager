FROM golang:1.17-alpine AS build
# For alpine image need to install bash
RUN apk add --no-cache bash
WORKDIR /go/src
COPY db /go/src/db
COPY go /go/src/go
COPY models /go/src/models
COPY util /go/src/util
COPY main.go /go/src
COPY go.mod /go/src
COPY subscription_manager.yaml /go/src

ENV CGO_ENABLED=0
RUN go get subscriptionManager
RUN go build -a -installsuffix cgo -o subscriptionManager .

EXPOSE 8085/tcp

# For troubleshooting uncomment this section and comment out the section production-grade version.
# Refer to the file startService.sh for what it actually starts.
ENV GIN_MODE=release
ADD startService.sh .
RUN chmod +x startService.sh
ENTRYPOINT ["./startService.sh"]

# Production-grade version
#FROM scratch AS runtime
#ENV GIN_MODE=release
#COPY --from=build /go/src/subscriptionManager ./
#COPY resources ./resources

#ENTRYPOINT ["./subscriptionManager"]