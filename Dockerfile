# syntax=docker/dockerfile:experimental

# From golang:1.16.3-buster

# COPY firestore.go firestore.go
# RUN go mod init "firebase.google.com/go/v4"
# RUN go mod download
# # RUN go get firebase.google.com/go/v4
# # RUN go get "google.golang.org/api/option"
# CMD go run firestore.go