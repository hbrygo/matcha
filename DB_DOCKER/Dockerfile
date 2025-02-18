# go 1.24 
FROM golang:1.24-alpine

# necessary packages
RUN apk add --no-cache \
    sqlite \
    sqlite-dev \
    gcc \
    musl-dev \
    git

# working directory
WORKDIR /app

# Copy files 
COPY go.mod go.sum ./
COPY . .

# install go dependancies
RUN go mod download

# creation DB folder
RUN mkdir -p /app/data

# construction of app 
RUN CGO_ENABLED=1 go build -o main .

# port 8181
EXPOSE 8181

CMD ["./main"]