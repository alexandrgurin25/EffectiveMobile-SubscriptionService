FROM golang:1.25.3-alpine

# Устанавливаем рабочую директорию
WORKDIR /subscriptions

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Переходим в директорию с main.go и собираем приложение
RUN cd cmd/app && go build -o subscriptions .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./cmd/app/subscriptions"]