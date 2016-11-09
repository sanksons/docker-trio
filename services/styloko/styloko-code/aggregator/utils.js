module.exports = {
    getUrls : function(product_url, data, ty, chunkSize){
        var prefix = "id"
        if(ty=="sku"){
            prefix="sku";
        }
        data=data.split(",");
        var urls = []
        if (chunkSize==undefined){
            chunkSize=10
        }
        for (var x=0; x<data.length; x+=chunkSize){
            sl = data.slice(x, x+chunkSize);
            urls.push(product_url+"?"+prefix+"=["+sl+"]")
        }
        return {urls:urls,skus:data}
    },
    mobApi: function(skus){
        keys = "product_"+skus.join(",product_")
        mobapi_url = styloko.mobapiUrl+"mobapi/main/getKeys/?keys="+keys
        var prices = {}
        var promise = new Promise(function(resolve,reject){
            HttpService.get(mobapi_url,styloko.mobapiHeaders).then((res)=>{
                var data = JSON.parse(res.data).metadata.data
                for (x in data){
                    sku=x.split("_")[1]
                    var baseObj = {}
                    for (y in data[x].simples){
                        var simple_sku = y
                        var tmp = data[x].simples[y].meta
                        var spObj = specialPrice()
                        if(tmp.price){
                            spObj.price=parsePrice(tmp.price);
                        }
                        if (tmp.special_price){
                            spObj.specialPrice=parsePrice(tmp.special_price);
                        }
                        if (tmp.special_from_date){
                            spObj.specialPriceFromDate=tmp.special_from_date+" 00:00:00";
                        }
                        if (tmp.special_to_date){
                            spObj.specialPriceToDate=tmp.special_to_date+" 23:59:58";
                        }
                        baseObj[y]=spObj
                    }
                    var meta = data[x].meta
                    var price_info = {
                        max_price: parsePrice(meta.max_price),
                        price: parsePrice(meta.price),
                        max_original_price:parsePrice(meta.max_original_price),
                        original_price: parsePrice(meta.original_price),
                        special_price: parsePrice(meta.special_price),
                        max_special_price:parsePrice(meta.max_special_price),
                        max_saving_percentage:custParseFloat(meta.max_saving_percentage),
                        special_price_from:baseObj[y].specialPriceFromDate,
                        special_price_to: baseObj[y].specialPriceToDate,
                    }
                    baseObj.price_info = price_info;
                    prices[sku]=baseObj
                }
                resolve(prices)
            },(err)=>{
                reject(err)
            })
        })
        return promise
    },
    replacePrices(original, prices){
        for (x in original.data){
            var tmp = original.data[x];
            if(prices[tmp.sku]){
                for (y in tmp.simples){
                    if (prices[tmp.sku][tmp.simples[y].sku]){
                        tmp.simples[y].special_price=prices[tmp.sku][tmp.simples[y].sku].specialPrice
                        tmp.simples[y].special_price_from=prices[tmp.sku][tmp.simples[y].sku].specialPriceFromDate
                        tmp.simples[y].special_price_to=prices[tmp.sku][tmp.simples[y].sku].specialPriceToDate
                        tmp.simples[y].meta.price=prices[tmp.sku][tmp.simples[y].sku].price
                        tmp.simples[y].meta.original_price=prices[tmp.sku][tmp.simples[y].sku].price
                    }else{
                        tmp.simples[y].special_price=null
                        tmp.simples[y].special_price_from=null
                        tmp.simples[y].special_price_to=null
                    }
                }
                tmp.price_info = prices[tmp.sku].price_info
            }else{
                for (y in tmp.simples){
                    tmp.simples[y].special_price=null
                    tmp.simples[y].special_price_from=null
                    tmp.simples[y].special_price_to=null
                }
                tmp.price_info.special_price=null
                tmp.price_info.max_special_price=null
                tmp.price_info.special_price_from=null
                tmp.price_info.special_price_to=null
                tmp.price_info.max_saving_percentage=0
            }
            original.data[x]=tmp
        }
    }
}

function specialPrice(){
    return {
        price:null,
        specialPrice:null,
        specialPriceToDate:null,
        specialPriceFromDate:null,
    }
}

function parsePrice(num){
    if (!isNaN(num) && num!=null && num!=undefined){
        return parseFloat(num)
    }
    return null
}

function custParseFloat(num){
    if (!isNaN(num) && num!=null && num!=undefined){
        return parseFloat(num)
    }
    return 0
}
