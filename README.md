### Миграции
`migrate -database postgresql://localhost:5432/archive?sslmode=disable"&"user=local"&"password=local_password -path migrations up`

`migrate create -ext sql -dir migrations %name%`