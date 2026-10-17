FROM golang:1.26.5 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/server ./cmd
RUN mkdir /out/storage

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build --chown=65532:65532 /out/storage /app/storage
COPY --from=build /out/server /app/server
USER 65532:65532
EXPOSE 8080 8081
ENTRYPOINT ["/app/server"]
