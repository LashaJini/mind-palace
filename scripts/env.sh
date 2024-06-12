#!/bin/bash

if [[ -z "${MP_ENV}" ]]; then
	echo "MP_ENV variable is not set. Exiting..."
	exit 1
fi

case $MP_ENV in
prod)
	source ./.env.prod
	;;
dev)
	source ./.env.dev
	;;
test)
	source ./.env.test
	;;
*)
	echo "Unknown MP_ENV value $MP_ENV. Exiting..."
	exit 1
	;;
esac
