# @jopit
# Utilizo esta imagen como builder para las dependencias privadas
FROM golang:alpine as builder

RUN apk --update add ca-certificates

RUN apk add build-base
WORKDIR /builder
ADD . /builder
WORKDIR /builder/environment
WORKDIR /builder/src/main/api
RUN CGO_ENABLED=0 go build -mod=vendor

# Corro la app aca ya que es la imagen mas liviana existente
FROM alpine
COPY --from=builder /builder/src/main/api/api /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/app"]