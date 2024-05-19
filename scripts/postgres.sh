#!/bin/bash

source ./.env

NAME=postgres13
POSTGRES_PASSWORD=$DB_PASS
POSTGRESQL_URL="postgres://$DB_USER:$POSTGRES_PASSWORD@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
MIGRATIONS_DIR="migrations"

start() {
	docker run --rm -d \
		-v $HOME/sql-db:/data \
		--name $NAME \
		-p $DB_PORT:$DB_PORT \
		-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
		postgres:13.14 1>/dev/null

	if [ $? -ne 0 ]; then
		echo "Start failed."
		exit 1
	fi
	echo "Success"
}

stop() {
	docker stop $NAME 1>/dev/null

	if [ $? -ne 0 ]; then
		echo "Stop failed."
		exit 1
	fi
	echo "Success"
}

cli() {
	docker exec -it postgres13 \
		bash -c "echo 'set -o vi'>~/.bashrc && \
              echo 'set editing-mode vi'>~/.inputrc && \
              su - postgres -c 'psql -d $DB_NAME'"
}

# run once, when application is created
first() {
	docker exec postgres13 \
		bash -c "su - postgres -c 'createdb $DB_NAME'"
}

db:migrate() {
	case $1 in
	create)
		db:migrate:create $2
		;;
	up)
		db:migrate:up $2
		;;
	down)
		db:migrate:down $2
		;;
	*)
		echo "please use bash postgres.sh db:migrate create|up|down"
		;;
	esac
}

db:migrate:create() {
	migrate create -ext sql -dir $MIGRATIONS_DIR $1 1>/dev/null

	if [ $? -ne 0 ]; then
		echo "db:migrate create failed."
		exit 1
	fi
	echo "Success"
}

db:migrate:up() {
	migrate -database $POSTGRESQL_URL -path $MIGRATIONS_DIR up

	if [ $? -ne 0 ]; then
		echo "db:migrate up failed."
		exit 1
	fi
	echo "Success"
}

db:migrate:down() {
	yes | migrate -database $POSTGRESQL_URL -path $MIGRATIONS_DIR down

	if [ $? -ne 0 ]; then
		echo "db:migrate down failed."
		exit 1
	fi
	echo "Success"
}

case $1 in
start)
	start
	;;
stop)
	stop
	;;
cli)
	cli
	;;
first)
	first
	;;
db:migrate)
	$@
	;;
*)
	echo "please use bash postgres.sh start|stop|cli|first|db:migrate"
	;;
esac
