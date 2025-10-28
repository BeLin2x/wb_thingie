# wb_thingie
больше для меня памятка

ЗАПУСК СЕРВЕРА (ВРУБИТЬ ДОКЕРА И ПОСТГРЕС)
  docker run -d --name nats-streaming -p 4222:4222 -p 8222:8222 nats-streaming:latest

  go run cmd/server/main.go

САЙТ
  http://localhost:8080/
  b563feb7b2b84b6test

ТЕСТЫ
  unit тест
    go test ./internal/cache/ -v
  стресс тест
    vegeta attack -duration=10s -rate=100 -targets=scripts/stress_test.txt | vegeta report
    vegeta attack -duration=5s -rate=500 -targets=scripts/stress_test.txt | vegeta report

СОХРАНИТЬ ТЕСТЫ
  vegeta attack -duration=10s -rate=200 -targets=scripts/stress_test.txt > results.bin
  vegeta report results.bin
  vegeta plot results.bin > plot.html
