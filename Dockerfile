FROM golang:alpine AS build
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go build ./cmd/cmtpr -o /bin/cmtpr

FROM alpine
COPY --from=build /bin/cmtpr /bin/cmtpr
ENTRYPOINT ["cmtpr"]