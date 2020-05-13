FROM golang:1.14 as builder
RUN mkdir /go/src/sshhoneypot
ADD . /go/src/sshhoneypot
WORKDIR /go/src/sshhoneypot
RUN CGO_ENABLED=0 GOOS=linux go build -o /sshhoneypot sshhoneypot.go


FROM scratch as final
WORKDIR /
COPY --from=builder /sshhoneypot /
EXPOSE 2222
ENTRYPOINT []
CMD ["/sshhoneypot"]