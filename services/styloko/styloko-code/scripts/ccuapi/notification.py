#!/usr/bin/python

import config
import requests
import json
from socket import gethostname

api_base_url = 'https://app.datadoghq.com'
event_url = '/api/v1/events'
series_url = '/api/v1/series'

def send_notification(url, ty):
    event = {
        "Title":ty.title()+ " cache invalidation failed",
        "Text":"Cache invalidation failure for url: "+url+"\n Alerting:\n"+config.datadog_conf["alert"],
        "Tags":[ty, "cache-invalidation-failure","styloko"],
        "AlertType":"error",
        "Host":gethostname(),
        }
    payload = {"api_key":config.datadog_conf["api_key"], "application_key":config.datadog_conf["app_key"]}
    api_url = api_base_url+event_url
    r=requests.post(api_url, params=payload, data=json.dumps(event))
    print(r.status_code, r.reason)
    return r

def test_notification():
    url="http://test.url"
    ty="varnish"

    send_notification(url, ty)

if __name__=="__main__":
    test_notification()
