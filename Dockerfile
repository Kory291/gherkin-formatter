FROM golang:latest

WORKDIR /workdir

COPY go.sum go.sum
COPY go.mod go.mod

COPY cmd cmd
COPY internal internal
COPY main.go main.go

RUN go build

ENTRYPOINT [ "./gherkin-formatter" ]
CMD [ "" ]