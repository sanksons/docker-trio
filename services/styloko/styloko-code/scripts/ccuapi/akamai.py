#!/usr/bin/python

############################################
# @author: Ishaan Bahal
# @email: ishaan.bahal@jabong.com
############################################

from ccuapi.purge import PurgeRequest
import json
import web
import sys
from purge_varnish import purge_varnish,url_builder
import urlparse
from config import akamai_conf, varnish_hosts
from notification import send_notification

# email can be an array of emails, must be registered on Akamai
def purgerInit(**kwargs):
    if "username" in kwargs and "password" in kwargs and "email" not in kwargs:
        purger = PurgeRequest(username=kwargs["username"], password=kwargs["password"], kind="arl")
        print("Purger init successful")
        return purger
    elif "username" in kwargs and "password" in kwargs and "email" in kwargs:
        purger = PurgeRequest(username=kwargs["username"], password=kwargs["password"], email=kwargs["email"])
        print("Purger init successful")
        return purger
    else:
        print("Username, password not provided to purger")
        sys.exit(1)

# serialize raises Exception for invalid JSON
def serialize(data):
    data = json.loads(data)
    if "objects" in data and isinstance(data["objects"], list):
        return data
    else:
        raise Exception("JSON invalid, must contain objects array")


routes = (
    '/purge', 'purge_urls',
    '/invalidate', 'invalidate_urls',
    '/varnish','varnish',
)

class varnish:
    def GET(self):
        '''
        Provide query param with URL.
        Possible params include url (complete with schema) or endpoint
        Examples:
        url=http://127.0.0.1:6081/getCategories
        endpoint=/getCategories
        '''
        user_data = web.input()

        if "url" in user_data:
            url=user_data.url
            return purge_varnish(url)
        elif "endpoint" in user_data:
            endpoint = user_data.endpoint
            responses = []
            for x in varnish_hosts:
                url = urlparse.urljoin(x, endpoint)
                responses.append(purge_varnish(url))
            return json.dumps(responses)
        else:
            return "Please provide a URL like: \n/varnish?url=www.example.com\n/varnish/?endpoint=catalog/v1/products/1"

    def POST(self):
        '''
        POST request requires data in the form:
        {
            "endpoints":["/getCategories"]
        }
        '''
        data =web.data()
        data = json.loads(data)
        if "endpoints" not in data:
            return "Please provide an array of endpoints"
        responses=[]
        for y in data["endpoints"]:
            for x in varnish_hosts:
                url = urlparse.urljoin(x, y)
                responses.append(purge_varnish(url))
        return json.dumps(responses)

class invalidate_urls:
    def GET(self):
        retData = '''
        Only POST api, please post array of URLs to invalidate.
        Must be of the form:

        {
            objects: [
                'http://www.example.com/',
                'www.example.com/static.png'
                ]
        }
        '''
        return retData
    def POST(self):
        # Core business logic will be here.
        return "Not implemented yet"

class purge_urls:
    def GET(self):
        retData = '''
        Only POST api, please post array of URLs to purge.
        Please note: PURGE requests take a lot of time, you may need invalidate only.
        To invalidate, POST data to endpoint: /invalidate
        Must be of the form:

        {
            objects: [
                'http://www.example.com/',
                'www.example.com/static.png'
                ]
        }
        '''
        return retData
    def POST(self):
        data =web.data()
        try:
            res = serialize(data)
            results=purge(res["objects"])
        except Exception as e:
            return e.message
        return "posted: "+ str(data)

app = web.application(routes, globals())

# username, password, varnish_hosts = configRead('conf.json')
if akamai_conf["username"]!="" and akamai_conf["password"]!="":
    purger=purgerInit(username=akamai_conf["username"], password=akamai_conf["password"])

def purge(keys):
    purger.add(keys)
    results = purger.purge()
    return results

def main():
    app.run()

if __name__ == "__main__":
    main()
