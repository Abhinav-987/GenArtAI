# Build Stage
FROM golang:1.22-alpine as builder

WORKDIR /app
RUN apk add --no-cache make nodejs npm git

COPY . ./
RUN make install
RUN make build

# Final Stage
FROM scratch
COPY --from=builder /app/genartai /genartai
COPY --from=builder /app/view /view
COPY --from=builder /app/public /public

EXPOSE 3000
ENTRYPOINT [ "/genartai" ]
