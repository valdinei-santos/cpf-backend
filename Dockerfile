#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:1.24.0-bullseye AS builder

# Create appuser.
ENV USER=backend
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o backend cmd/api/main.go

#####################################
#   STEP 2 build a small image      #
#####################################
FROM alpine:latest
ENV TZ="America/Sao_Paulo"
WORKDIR /app
COPY --from=builder /app/backend .
#COPY --from=builder /app/templates /app/templates
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use an unprivileged user.
USER backend:backend

CMD ["./backend"]  
EXPOSE 8800