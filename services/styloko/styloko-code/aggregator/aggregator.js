// aggregator.js
// Styloko to Boutique response mapper

var express = require("express");
var mapper = require("./mapper.js");
var aggregator = express();
var colors = require("colors");
var compression = require('compression')
var HttpService = require("./HttpService.js");

global.styloko = {};
global.styloko.baseUrl = process.env.PRODUCT_URL || "http://127.0.0.1:8084/"
global.styloko.mobapiUrl = process.env.MOBAPI_URL || "http://mobapi.jabong.com/"
global.styloko.mobapiHeaders = {
    "X-ROCKET-MOBAPI-TOKEN": "b5467071c82d1e0d88a46f6e057dbb88",
    "X-ROCKET-MOBAPI-VERSION": "3.0",
    "X-ROCKET-MOBAPI-PLATFORM": "application/ios.rocket.SHOP_MOBILEAPI_STAGING-v1.0+json",
    "X-USER-DEVICE-TYPE": "wsoa",
}
var chunksize = process.env.MAPPER_CHUNK_SIZE || 10
global.styloko.chunksize = parseInt(chunksize)
global.HttpService = HttpService;

try{
    redis_url = process.env.REDIS_STOCK_HOSTS || "redis://127.0.0.1:6379"
    var Redis = require('ioredis');
    if (redis_url.split(',').length==1){
            redis_ttl = process.env.REDIS_TTL || 604800;
            client = new Redis(redis_url);
            client.ttl(redis_ttl);
            global.styloko.redis = client
        }
        else{
            var nodes = [];
            var hosts = redis_url.split(",");
            for (var x in hosts){
                var tmp = hosts[x].split(":");
                nodes.push({host:tmp[0],port:tmp[1]})
            }
            redis_ttl = process.env.REDIS_TTL || 604800
            client = new Redis.Cluster(nodes,
                {
                    redisOptions: {
                        password: process.env.REDIS_PASSWORD||"",
                    }
                }
            );
            client.ttl(redis_ttl);
            global.styloko.redis = client
        }
    }
catch(e){
    console.log(e.message)
    console.log("Please install ioredis using: npm install ioredis --save");
    process.exit(1);
}
aggregator.use(compression());
aggregator.get('/catalog/v1/products/multi/:ids', function(req, res){
    mapper.map(req, res, "id",false);
});

aggregator.get('/catalog/v1/products/multi/sku/:skus', function(req, res){
    mapper.map(req, res, "sku",false);
})


aggregator.get('/catalog/v1/products/sku/:skus', function(req, res){
    mapper.map(req, res, "sku", true);
});

aggregator.get('/catalog/v1/products/:ids', function(req, res){
    mapper.map(req, res, "id", true);
});

var port = process.env.AGGREGATOR_PORT || 8088
console.log(("Started aggregator service at port: "+ port).green)
aggregator.listen(port);
