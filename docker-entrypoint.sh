#!/bin/sh

_main() {
	# if first arg looks like a flag, assume we want to run server server
	if [ "${1:0:1}" = '-' ]; then
		set -- go-http-server "$@"
	fi

	exec "$@"
}

_main "$@"