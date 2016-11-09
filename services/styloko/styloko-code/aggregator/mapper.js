var Promise = require("promise");
var transform = require("./transform.js");
var utils = require("./utils.js")
module.exports = {
    map: function (request, response, ty, flag) {
        var product_url = styloko.baseUrl + "catalog/v1/products/"
        var sizechart_url = styloko.baseUrl + "catalog/v1/sizechart/"
        var predef = false
        var mobapi_enable=false
        if (request.headers["visibility-type"]|| request.headers["Visibility-Type"]){
            predef =true
        }
        if (request.headers["mobapi-enable"]){
            if (request.headers["mobapi-enable"].toLowerCase()=="true"){
                mobapi_enable=true
            }
        }
        if (!predef){
            request.headers["Visibility-Type"] = "MULTI-SKU"
        }
        request.headers["Expanse"] = "XLarge"
        var small = false
        var catalog = false
        if (flag && !predef){
            request.headers["Visibility-Type"] = "PDP"
        }
        if (request.query.small === "true") {
            small = true;
        } else if (request.query.catalog === "true") {
            catalog = true;
        }
        var chunksize=styloko.chunksize
        if (request.query.chunksize){
            chunksize=parseInt(request.query.chunksize)
        }
        var stock = true
        if (request.query.stock){
            if (request.query.stock.toLowerCase()=="false"){
                stock=false
            }
        }
        var parentPromiseArray = [];
        var mobApiData = {};
        switch (ty) {
        case "id":
            this.urls = utils.getUrls(product_url, request.params.ids, "ids", chunksize).urls;
            break
        case "sku":
            var urls = utils.getUrls(product_url, request.params.skus, "sku", chunksize);
            this.urls = urls.urls
            if (mobapi_enable){
                var mobPro = utils.mobApi(urls.skus);
                mobPro.then((data)=>{
                    mobApiData = data
                },(err)=>{
                    mobApiData = {}
                })
                parentPromiseArray.push(mobPro);
            }
            break
        default:
            return response.status(404).send({
                "error": "Invalid URL"
            })
        }
        if (small) {
            request.headers["Expanse"] = "Large";
        } else if (catalog) {
            request.headers["Expanse"] = "Catalog"
        }
        var btRes = {
            data: []
        };
        for (x in this.urls) {
            var promiseArray = [];
            var pro = HttpService.get(this.urls[x], request.headers);
            tmp_pro = new Promise(callback);
            parentPromiseArray.push(tmp_pro);
            function callback(resolve, reject){
                pro.then((res) => {
                    if (small) {
                        var td = transform.transformDataSmall(res.data)
                    } else if (catalog) {
                        var td = transform.transformDataCatalog(res.data)
                    } else {
                        var td = transform.transformData(res.data)
                    }
                    if (stock){
                        for (i in td.data) {
                            var pipeline = styloko.redis.pipeline();
                            for (j in td.data[i].simples) {
                                var stock_key = "{stock}_" + td.data[i].simples[j].id;
                                pipeline.hgetall(stock_key);
                            }
                            var p2 = pipeline.exec();
                            promiseArray.push(p2);
                            (function (i) {
                                p2.then(function (res) {
                                    for (j in res) {
                                        if (!res[j].length > 0) {
                                            td.data[i].simples[j].quantity = 0;
                                            continue;
                                        }
                                        var quantity = 0
                                        if (res[j][1].quantity != undefined) {
                                            quantity = res[j][1].quantity;
                                        }
                                        var reserved = 0
                                        if (res[j][1].reserved != undefined) {
                                            reserved = res[j][1].reserved;
                                        }
                                        var diff = quantity - reserved;
                                        if (diff > 0) {
                                            td.data[i].simples[j].quantity = diff;
                                        } else {
                                            td.data[i].simples[j].quantity = 0;
                                        }
                                    }
                                }, function (err) {
                                    console.log("Promise rejected by Redis")
                                    td.data[i].simples[j].quantity = 0;
                                })
                            }(i));
                            if (catalog) {
                                for (j in td.data[i].product_group) {
                                    var pipeline2 = styloko.redis.pipeline();
                                    for (x in td.data[i].product_group[j].simples) {
                                        var stock_key = "{stock}_" + td.data[i].product_group[j].simples[x].id;
                                        pipeline2.hgetall(stock_key);
                                    }
                                    var p3 = pipeline2.exec();
                                    promiseArray.push(p3);
                                    var pg_qty = [];
                                    (function (i, j) {
                                        p3.then((res) => {
                                            for (x in res) {
                                                if (!res[x].length > 0) {
                                                    td.data[i].product_group[j].simples[x].quantity = 0;
                                                    continue;
                                                }
                                                var quantity = 0
                                                if (res[x][1].quantity != undefined) {
                                                    quantity = res[x][1].quantity;
                                                }
                                                var reserved = 0
                                                if (res[x][1].reserved != undefined) {
                                                    reserved = res[x][1].reserved
                                                }
                                                var diff = quantity - reserved;
                                                if (diff > 0) {
                                                    td.data[i].product_group[j].simples[x].quantity = diff;
                                                } else {
                                                    td.data[i].product_group[j].simples[x].quantity = 0;
                                                }
                                            }
                                        }, (err) => {
                                            console.log(err);
                                        })
                                    }(i, j))
                                }
                            }
                        }
                    }

                    if (td.err) {
                        console.log({
                            success: false,
                            error_code: 500,
                            message: td.err,
                        })
                        return response.status(500).send({
                            success: false,
                            error_code: 500,
                            message: td.err,
                        })
                        reject();
                    }
                    Promise.all(promiseArray).then((data) => {
                    if (stock){
                        for (i in td.data) {
                            var qty = 0;
                            if (td.data[i].simples != undefined) {
                                for (j in td.data[i].simples) {
                                    qty += td.data[i].simples[j].quantity
                                }
                                td.data[i].quantity = qty;
                            }
                            pg_qty = 0;
                            for (j in td.data[i].product_group) {
                                if (td.data[i].product_group[j].simples != undefined) {
                                    for (x in td.data[i].product_group[j].simples) {
                                        pg_qty += td.data[i].product_group[j].simples[x].quantity;
                                    }
                                }
                                td.data[i].product_group[j].quantity = pg_qty;
                            }
                        }
                    }
                    btRes.data.push.apply(btRes.data,td.data);
                    resolve();
                })
                }, (err) => {
                    return response.negotiate(err)
                    reject();
                })
            }
        }
        Promise.all(parentPromiseArray.map(p => p.catch(e => e))).then((data) => {
            if (btRes.data.length == 0) {
                console.log({
                    success: false,
                    error_code: 100,
                    sku:"",
                    reason:"no data",
                    message: "Resource not found for given parameter",
                })
                return response.json({
                    success: false,
                    error_code: 100,
                    message: "Resource not found for given parameter",
                })
            }

            if(mobapi_enable && ty=="sku"){
                utils.replacePrices(btRes,mobApiData)
            }

            btRes.success = true;
            btRes.message = "";
            btRes.total = btRes.data.length;
            if (flag == true) {
                btRes.data = btRes.data[0]
                if (!btRes.data.visibility){
                    console.log({
                        success: false,
                        error_code: 100,
                        sku:btRes.data.sku,
                        message: "Product not visible",
                    })
                    return response.json({
                        success: false,
                        error_code: 100,
                        message: "Product not visible",
                    })
                }
            }
            return response.json(btRes)
        }, (err) => {
            return response.json(btRes)
        }).catch(e => console.log(e));
    }
}
