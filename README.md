# golang-mq

Basic of creating messaging application based on [RabbitMQ documentation](https://www.rabbitmq.com/getstarted.html) with default configuration.
This project using docker to run, so you have to installed it first

## Configuration

- Clone this repository <br> `git clone https://github.com/YogiPristiawan/golang-mq`<br><br>
- Install needed dependencies <br> `go mod tidy`<br><br>
- Run the application with docker compose <br> `docker compose up -d`<br><br>
- Make sure the application running on port `5672` and `15672` for admin management
