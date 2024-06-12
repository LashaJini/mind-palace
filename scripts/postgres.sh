#!/bin/bash

source ./.env

# TODO: this sucks
MIND_PALACE_USER=$(cat $HOME/.mind-palace/.info.json | jq -r '.current_user')
if [ $? -eq 1 ]; then
	echo "Can't get user."
	exit 1
fi

NAME=postgres13
POSTGRES_PASSWORD=$DB_PASS
DB_NAME="${MIND_PALACE_USER}${DB_NAME}"
POSTGRESQL_URL="postgres://$DB_USER:$POSTGRES_PASSWORD@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
MIGRATIONS_DIR="migrations"
VERSION=$DB_VERSION

start() {
	docker run --rm -d \
		-v $HOME/sql-db:/data \
		--name $NAME \
		-p $DB_PORT:$DB_PORT \
		-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
		postgres:$VERSION 1>/dev/null

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
	             psql -U $DB_USER -d $DB_NAME"
}

# run once, when application is created
first() {
	docker exec postgres13 \
		bash -c "su - $DB_USER -c 'createdb $DB_NAME'"
}

db:migrate() {
	mkdir $MIGRATIONS_DIR -p

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
	version)
		db:migrate:version
		;;
	fix)
		db:migrate:fix $2
		;;
	*)
		echo "please use bash postgres.sh db:migrate create|up|down|version|fix"
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
	migrate -database $POSTGRESQL_URL -path $MIGRATIONS_DIR down 1

	if [ $? -ne 0 ]; then
		echo "db:migrate down failed."
		exit 1
	fi
	echo "Success"
}

db:migrate:version() {
	migrate -database $POSTGRESQL_URL -path $MIGRATIONS_DIR version

	if [ $? -ne 0 ]; then
		echo "db:migrate version failed."
		exit 1
	fi
	echo "Success"
}

db:migrate:fix() {
	migrate -database $POSTGRESQL_URL -path $MIGRATIONS_DIR force $1

	if [ $? -ne 0 ]; then
		echo "db:migrate fix failed."
		exit 1
	fi
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
