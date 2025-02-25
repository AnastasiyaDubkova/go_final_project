FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

# При необходимости можно использовать пер.окружения, они перезатрут данные из файла .env
ENV TODO_PORT=
ENV TODO_DBFILE=

# НЕ будем использовать конструкцию EXPOSE, тк пор может изменяться из пер.окружения и/или .env
# вместо этого опишем запуск команды `docker run` в README.md
# EXPOSE 7540

CMD ["/my_app"]