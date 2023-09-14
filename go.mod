module github.com/1jack80/todo-api

go 1.21.0

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-sql-driver/mysql v1.7.1
)

require github.com/1jack80/guardian v0.1.4

replace github.com/1jack80/guardian => ../../golang/guardian
