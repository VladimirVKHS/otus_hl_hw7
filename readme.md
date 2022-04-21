# Otus Highload Architect homework 7

Websocket server for social network (homework 1)

Взаимодействие основного приложения с websocket-сервером реализуется через RabbitMQ.
Компонент может быть горизонтально масштабирован.
На экземпляр компонента поступают сообщения только для тех пользователей, которые подключены к нему. Это достигается за счет применения отдельной очереди на каждый экземпляр, связанной с общим exchange индивидуальным набором байндингов (routing key по id пользователя)

## Run

    docker-compose up

## .env configuration

    JWT_SECRET={any_string}}
    WS_PORT=8086
    RABBIT_HOST=otus_sn_rabbitmq
    RABBIT_PORT=5672
    RABBIT_WS_EXCHANGE=ws
