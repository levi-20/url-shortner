FROM golang:1.25.10-alpine3.23 AS base

WORKDIR /app

COPY go.mod go.sum main.go .

RUN go mod download

# https://pkg.go.dev/cmd/link
# Removing the symbol tables reduces the size of the final binary by about ~5MB
# -s Omit the symbol table and debug information
# -w Omit the DWARF symbol table
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./main main.go

# Reset the context to "scratch" to create a super tiny docker image
FROM scratch

WORKDIR /app
COPY ./db ./db
COPY --from=base /app/main ./

# The final binary is still about 7MB in a docker image that CONTAINS NOTHING
CMD [ "./main" ]
