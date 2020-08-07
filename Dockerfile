FROM alpine:3.11.6 AS build
RUN apk add --no-cache go
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o app pix-api

FROM alpine:3.11.6
WORKDIR /app
COPY --from=build /app/app /app/app
CMD ["/app/app"]
