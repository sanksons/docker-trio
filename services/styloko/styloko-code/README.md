## Styloko Readme 

### Requirements
- MongoDB Local installation 
- MySql Local installation
- Datadog agent (not mandatory)
- GoLang (**latest version**)
- Redis-Server Local installation


### Setup
- Clone styloko in your `$GOPATH/src/github.com/jabong/` using: `git clone git@github.com:jabong/styloko.git`

- Go to the git directory just cloned, and type `make initdev`. This will create these files in your setup:
```
config/newApp/dev.json
config/newApp/api_dev.json
config/logger/logger_dev.json
scripts/env_dev.sh
```

- Edit `scripts/env_dev.sh` and `config/newApp/dev.json` to your local settings.

- Add these line to your `~./bashrc` : 
    - `export FLORESTENV="DEV"`
    - `source $GOPATH/src/github.com/jabong/src/styloko/scripts/env_dev.sh`

    > NOTE: source command path should change depending on your installation directory.

- Type this command: `source ~/.bashrc`

- Type the following command in your terminal
    `mkdir -p /var/log/styloko/;sudo chown -R $USER:$USER /var/log/styloko/;`

- Hoping every setting is correct, run command `make run` to run the project.

- Visit http://localhost:8084/catalog/healthcheck to see if the service is up and running.

- **[Optional]** Run migrations. For help in running migrations, see below.

### Running Migrations

> 1. Migrations are a big task and might take a lot of time.
> 2. Migrations are required by styloko since it uses Mongo as its primary database.

***API based migrations***

**Endpoint:** `catalog/v1/migrations/?key=categories`

**Possible keys**

- brands
- filters
- categories
- categorysegments
- attributesets
- attributes
- products
- productgroups
- productsindex
- productsdrop
- productsactive
- productsinactive
- taxclass
- sizecharts
- productsizechart

**Endpoint with ID values:**
`catalog/v1/migrations/?key=productsbyid&id=35000`

**Keys with ID are**

- productsbyseller
- productsbyid
- attributesbyid

***CLI commands (Same keys as above)***

**With binary**

```
./styloko -m attributes
./styloko -m productsbyid -i 35000
./styloko -m="attributes categories brands sizecharts filters"
```

**With make run**

```
make a="-m attributes" run
make a="-m productsbyid -i 35000" run
make a="-m=\"attributes categories brands sizecharts filters\"" run
```

### Supported API Endpoints
**Base URL:** `http://host:port/catalog/v1/`

| Endpoints | Methods | Wiki Link |
|:----------|:--------|:----------|
|*brands*|`GET`, `PUT`, `POST`|http://wiki.jira.rocket-internet.de/display/INDFAS/Brands |
|*categories*|`GET`, `PUT`, `POST`|http://wiki.jira.rocket-internet.de/display/INDFAS/Category |
|*attributes*|`GET`, `PUT`, `POST`|http://wiki.jira.rocket-internet.de/display/INDFAS/Attribute |
|*attributeSets*|`GET`, `PUT`, `POST`|N/A|
|*bootstrap*|`POST`|N/A|
|*product*|`GET`, `PUT`, `POST`, `DELETE`|http://wiki.jira.rocket-internet.de/display/INDFAS/Product |
|*sizechart*|`POST`|N/A|
|*migrations*|`GET`|N/A|
|*catalogty*|`GET`|N/A|
|*standardsize*|`POST`|http://wiki.jira.rocket-internet.de/display/INDFAS/Standard+Size |
|*taxclass*|`GET`|N/A|

**For more details on the project visit:**

http://wiki.jira.rocket-internet.de/display/INDFAS/Styloko


### Cache flush APIs for Varnish and Akamai
Styloko provides varnish and Akamai flush APIs for provided Endpoints

**Steps to install**
- Run `pip install -f requirements.txt`, may require *sudo* if outside virtualenv
- Modify `config.py` to reflect the akamai settings and varnish list of hosts
- Varnish flush may not work directly, this code may need to be added to the `/etc/varnish/default.vcl` (**4.1 and above**):

    ```
    sub vcl_recv {
            if (req.method == "PURGE") {
                    if (!client.ip ~ purge) {
                            return(synth(405,"Not allowed."));
                    }
                    return (purge);
            }
    }

    ```

- Start the python server by typing either of these commands: 

```
python akamai.py
python akamai.py 1234
python akamai.py 127.0.0.1 1234
```
