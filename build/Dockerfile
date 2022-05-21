FROM golang:1.18-alpine3.15 as builder
LABEL description="mawinter-discord server"

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY go.mod /app/go.mod
RUN cd /app && go mod tidy
WORKDIR /app/cmd/mawinter-discord/
RUN go build -o /app/cmd/mawinter-discord/mawinter-discord

FROM alpine:3.15
COPY --from=builder /app/cmd/mawinter-discord/mawinter-discord /bin/mawinter-discord
ENTRYPOINT ["/bin/mawinter-discord", "run"]
