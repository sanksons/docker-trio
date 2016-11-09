#!/usr/bin/python
import os

__vhosts__ = os.getenv("VARNISH_HOSTS","http://127.0.0.1:6081/").split(",")

varnish_hosts = __vhosts__

akamai_conf = {
    "username":os.getenv("AKAMAI_USERNAME", ""),
    "password":os.getenv("AKAMAI_PASSWORD", ""),
    "email":os.getenv("AKAMAI_EMAIL", ""),
    }

datadog_conf = {
    "api_key":os.getenv("DATADOG_API_KEY","239cb79c7f1afddb0f60e7363ce0c53a"),
    "app_key":os.getenv("DATADOG_APP_KEY","735d8dbb85bc4b38eb9486c2ef9812f2c5a38e36"),
    "alert":os.getenv("NOTIF_ADDR","@apoorva.moghey@jabong.com")
}
