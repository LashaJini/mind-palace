#!/bin/bash

source ./scripts/env.sh

CONTAINER_NAME=postgres13
POSTGRES_PASSWORD=$DB_PASS
# POSTGRESQL_URL="postgres://$DB_USER:$POSTGRES_PASSWORD@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
VERSION=$DB_VERSION

start() {
	mkdir $MIGRATIONS_DIR $POSTGRES_DATA_DIR -p

	# https://hub.docker.com/_/postgres
	docker run --rm -d \
		--name $CONTAINER_NAME \
		-p $DB_PORT:$DB_PORT \
		-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
		-e POSTGRES_DB=$DB_NAME \
		-e POSTGRES_USER=$DB_USER \
		-v ./$POSTGRES_DATA_DIR:/var/lib/postgresql/data \
		postgres:$VERSION 1>/dev/null

	if [ $? -ne 0 ]; then
		echo "Start failed."
		exit 1
	fi

	echo "Success"
}

stop() {
	docker stop $CONTAINER_NAME 1>/dev/null

	if [ $? -ne 0 ]; then
		echo "Stop failed."
		exit 1
	fi
	echo "Success"
}

cli() {
	docker exec -it $CONTAINER_NAME \
		bash -c "echo 'set -o vi'>~/.bashrc && \
	             echo 'set editing-mode vi'>~/.inputrc && \
	             psql -U $DB_USER -d $DB_NAME"
}

drop() {
	docker exec $CONTAINER_NAME \
		bash -c "psql -U $DB_USER -d postgres -c 'drop database $DB_NAME'" &&
		sudo rm -rf $POSTGRES_DATA_DIR
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
drop)
	drop
	;;
*)
	echo "please use bash postgres.sh start|stop|cli|first|drop"
	;;
esac
