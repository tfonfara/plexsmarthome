#!/bin/sh

if [ -z ${GID} ]; then GID=1000; fi
if [ -z ${UID} ]; then UID=1000; fi

if [ "${GID}" -ne "0" ]; then
	GROUP=plex
	if [ -z $(getent group $GROUP) ]; then
	    addgroup -g $GID $GROUP
	fi
else
	GROUP=root
fi

if [ "$UID" -ne "0" ]; then
    USER=plex
    if [ -z $(getent passwd $USER) ]; then
        adduser -D -g '' -G ${GROUP} -u $UID ${USER} -s /bin/false
    fi
else
    USER=root
fi

chown -R ${UID}:${GID} ${PLEX_CONFIG}

su-exec ${UID}:${GID} "$@"
