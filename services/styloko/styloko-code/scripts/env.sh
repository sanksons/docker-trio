#!/usr/bin/env bash

#!/bin/bash

export MASTER_USERNAME="bob_live"
export MASTER_PASSWORD="AaegiQdv"
export MASTER_HOST="192.168.89.229"
export MASTER_DBNAME="bob_live"
export MASTER_MAX_OPEN_CON=100
export MASTER_MAX_IDLE_CON=5
export SLAVE_USERNAME="bob_live"
export SLAVE_PASSWORD="AaegiQdv"
export SLAVE_HOST="192.168.89.229"
export SLAVE_DBNAME="bob_live"
export SLAVE_MAX_OPEN_CON="100"
export SLAVE_MAX_IDLE_CON="5"
export MONGO_CONNECTION_URL="mongodb://styloko:KNlZmiaNUp0B@192.168.89.223:27017,192.168.89.223:27018"
export MONGO_DBNAME="styloko"
export REDIS_HOSTS="localhost:6379"
export REDIS_PASSWORD="123"
export REDIS_POOLSIZE="25"
export JBUS_URL="http://jabongbus:8087/omnibus"
export JBUS_PUB="catalog"
export JBUS_RKEY="#.erp.#"
export BOOTPOOL_WORKER_COUNT="100"
export BOOTPOOL_QUEUE_SIZE="1000"
export NOTIF_ADDR="@apoorva.moghey@jabong.com"
export CCAPI_USERNAME=""
export CCAPI_PASSWORD=""
export CCAPI_URL="http://ccapi"
export PRO_MIG_LIMIT=100
export DATADOG_ENABLE=false
export PRODUCT_SELLER_UPDATE=false
export SELLER_SKU_LIMIT=100
