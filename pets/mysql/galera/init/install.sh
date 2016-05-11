#! /bin/bash

# This volume is assumed to exist and is shared with parent of the init
# container. It contains the mysq config. In the future it
CONFIG_VOLUME="/etc/mysql"

# This volume is assumed to exist and is shared with the peer-finder
# init container. It contains on-start/change configuration scripts.
WORKDIR_VOLUME="/work-dir"

for i in "$@"
do
case $i in
    -c=*|--config=*)
    CONFIG_VOLUME="${i#*=}"
    shift
    ;;
    -w=*|--work-dir=*)
    WORKDIR_VOLUME="${i#*=}"
    shift
    ;;
    *)
    # unknown option
    ;;
esac
done

echo installing config scripts into "${WORKDIR_VOLUME}"
mkdir -p "${WORKDIR_VOLUME}"
cp /on-start.sh "${WORKDIR_VOLUME}"/
cp /peer-finder "${WORKDIR_VOLUME}"/

echo installing my-galera.cnf into "${CONFIG_VOLUME}"
mkdir -p "${CONFIG_VOLUME}"
chown -R mysql:mysql "${CONFIG_VOLUME}"
cp /my-galera.cnf "${CONFIG_VOLUME}"/
