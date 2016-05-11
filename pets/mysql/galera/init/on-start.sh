#! /bin/bash

# This script writes out a mysql galera config using a list of newline seperated
# peer DNS names it accepts through stdin.

# /etc/mysql is assumed to be a shared volume so we can modify my.cnf as required
# to keep the config up to date, without wrapping mysqld in a custom pid1.
# The config location is intentionally not /etc/mysql/my.cnf because the
# standard base image clobbers that location.
CFG=/etc/mysql/my-galera.cnf

function join {
    local IFS="$1"; shift; echo "$*";
}

HOSTNAME=$(hostname)
# Parse out cluster name, formatted as: petset_name-index
IFS='-' read -ra ADDR <<< "$(hostname)"
CLUSTER_NAME="${ADDR[0]}"

while read -ra LINE; do
    if [[ "${LINE}" == *"${HOSTNAME}"* ]]; then
        MY_NAME=$LINE
    fi
    PEERS=("${PEERS[@]}" $LINE)
done

if [ "${#PEERS[@]}" = 1 ]; then
    WSREP_CLUSTER_ADDRESS=""
else
    WSREP_CLUSTER_ADDRESS=$(join , "${PEERS[@]}")
fi
sed -i -e "s|^wsrep_node_address=.*$|wsrep_node_address=${MY_NAME}|" ${CFG}
sed -i -e "s|^wsrep_cluster_name=.*$|wsrep_cluster_name=${CLUSTER_NAME}|" ${CFG}
sed -i -e "s|^wsrep_cluster_address=.*$|wsrep_cluster_address=gcomm://${WSREP_CLUSTER_ADDRESS}|" ${CFG}

# don't need a restart, we're just writing the conf in case there's an
# unexpected restart on the node.
