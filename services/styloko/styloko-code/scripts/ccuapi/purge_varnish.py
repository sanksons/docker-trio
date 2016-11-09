#!/usr/bin/python

############################################
# @author: Ishaan Bahal
# @email: ishaan.bahal@jabong.com
############################################

import requests
from os import path
from notification import send_notification

def purge_varnish(url):
    '''Purge fires a request to provided URL with the method PURGE.
    PURGE is not a standard HTTP method

    url: URL on which a PURGE request will be fired
    '''
    resp = requests.request("PURGE",url)
    if (resp.status_code/10)==20:
         return {"status":True,"response":"success"}
    else:
        send_notification(url, "varnish")
        return {"status":False, "response":resp.reason}

def url_builder(ty, hostname, endpoint, **kwargs):
    '''
    Builds url based on ty provided
    hostname: hostname of the service being hit
    endpoint: endpoint of the service
    :: kwargs : q=: string with & seperated query values
    :: kwards : path=: string with extra path params after endpoint

    returns a URL for PURGE request
    '''
    path=""
    q=""
    base_url=""
    if "q" in kwargs:
        q="?"+kwargs["q"]
    if "path" in kwargs:
        path=kwargs["path"]

    # Ty cases below
    if ty=="styloko":
        base_url = path.join("http://",hostname,"catalog/v1",endpoint, path,q)
    return base_url
