FROM golang:1.25.10-alpine3.23 AS base

RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates

WORKDIR /app

# Since dependencies don't change a lot, we first just import go.mod and go.sum
COPY go.mod go.sum .

RUN go mod download

COPY . .

# https://pkg.go.dev/cmd/link
# Removing the symbol tables reduces the size of the final binary by about ~5MB
# -s Omit the symbol table and debug information
# -w Omit the DWARF symbol table
ARG COMMIT_SHA=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X github.com/levi-20/url-shortner/templates.CommitSHA=${COMMIT_SHA}" -o ./main main.go

FROM golang:1.25.10-alpine3.23

WORKDIR /app

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /app/main ./

# The final binary is still about 7MB in a docker image that CONTAINS NOTHING
# Golang is terrible when it comes to creating small binaries
CMD [ "./main" ]