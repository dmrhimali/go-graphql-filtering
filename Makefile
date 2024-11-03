include .env

#call: make create_migration ARG=create_author_friend_author_table
create_migration:
	~/go/bin/migrate	create	 -dir internal/db/migrations	-ext sql	-seq $(ARG)

migrate_up:
	~/go/bin/migrate	-path=internal/db/migrations	-database	"mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}"	-verbose	up

migrate_down:
	~/go/bin/migrate -path=internal/db/migrations -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}?"	-verbose	down

.PHONY: create_migration migrate_up migrate_down