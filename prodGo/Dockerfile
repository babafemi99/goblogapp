 FROM golang:latest

 LABEL maintainer='femi <ooluwa27@gmail.com>'

 WORKDIR /app

 COPY go.mod .

 COPY go.sum .

 RUN go mod download

 COPY . .

 ENV PORT=8100

 RUN go build

#  CMD [ "/appz" ]