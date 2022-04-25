FROM golang:1.18.1-alpine3.15 as builder
ADD . /app
WORKDIR /app/go-server/cmd/svr
RUN apk add git
RUN apk add build-base
RUN go mod download
RUN go build -o /main .

FROM node:16.14-buster AS node_builder
WORKDIR /code
COPY --from=builder /app .
RUN rm -rf /go-server/
RUN yarn install --frozen--lockfile
COPY . .
RUN yarn build

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
COPY --from=builder /main .
COPY --from=builder /app/go-server/cmd/svr/*config.json .
COPY --from=node_builder /code/build ./web
RUN chmod +x ./main
EXPOSE 80
CMD ./main
