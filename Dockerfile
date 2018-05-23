FROM golang:alpine
RUN mkdir -p /go/src/showhook/
WORKDIR /go/src/showhook/
ADD main.go .
RUN apk --no-cache add git
RUN go get -d -v
RUN go install -v
RUN go build main.go

FROM alpine:3.7
RUN mkdir -p /go/src/showhook/
WORKDIR /go/src/showhook/
COPY --from=0 /go/src/showhook/main .
ADD public public
EXPOSE 9000
ARG COMMIT
ENV COMMIT ${COMMIT}
CMD ["./main"]
