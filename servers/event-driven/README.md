# gRPC-сервис для создания опросов и голосования с лидербордом в реальном времени

## Технологии
- Go
- gRPC + Protocol Buffers
- PostgreSQL
- Redis
- Docker + docker-compose

# Запуск
```
docker-compose up -d
go run cmd/main.go
```

# Примеры запросов
```
# Создать опрос
grpcurl -plaintext -d "{\"title\":\"Вопрос\",\"description\":\"Описание\"}" localhost:9090 poll.PollService/CreatePoll

# Проголосовать
grpcurl -plaintext -d "{\"poll_id\":\"ID\",\"option_id\":\"ID\",\"user_id\":\"user1\"}" localhost:9090 vote.VoteService/CastVote

# Подписаться на лидерборд
grpcurl -plaintext -d "{\"poll_id\":\"ID\",\"top_n\":5}" localhost:9090 vote.VoteService/SubscribeLeaderboard
```