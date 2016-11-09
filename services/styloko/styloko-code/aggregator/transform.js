wh = {
    "AGR01": 49,
    "AMD01": 43,
    "AMD02": 50,
    "Bangalore PC": 12,
    "BLR01": 41,
    "BLR02": 40,
    "BLR03": 42,
    "BLR04": 51,
    "BLR05": 52,
    "BOMBAY PC": 4,
    "BRD01": 17,
    "CCU01": 29,
    "CHD01": 39,
    "CJB01": 48,
    "DEL01": 23,
    "Delhi PC": 3,
    "DELHIWH": 15,
    "Dropship": 5,
    "DWARKA WH": 11,
    "GGN01": 35,
    "GGNPC": 16,
    "GOI01": 36,
    "HYD01": 33,
    "HYD02": 19,
    "IDR01": 32,
    "JAI01": 31,
    "JUC01": 30,
    "Khawaspur": 13,
    "LUH01": 28,
    "LUK01": 46,
    "MAA01": 37,
    "MRT01": 45,
    "MUM01": 26,
    "MUM02": 25,
    "MUM03": 24,
    "NDA01": 22,
    "NMB01": 56,
    "PAT01": 55,
    "PNQ01": 21,
    "PNQ02": 20,
    "STV01": 18,
    "TPJ01": 38,
    "TUP01": 53,
    "UP01": 54,
    "Warehouse": 2,
    "XGGN01": 14,
}

var getAttrId = function (name) {
    return wh[name]
}

var newBoutiqueObj = function () {
    return {
        id: 0,
        sku: "",
        name: "",
        status: "",
        pet_approved: 0,

        attribute_set_name: {
            id: 0,
            name: "",
            label: "",
            label_en: null,
            identifier: "",
        },
        brand: {},
        catalog_ty: null,
        categories: [],
        product_group: null,
        sizechart: {
            id: null,
            type: null,
        },
        reward_points: 0,
        rating_info: {
            total: 0,
            single: "",
        },
        price_info: {
            max_price: 0,
            price: 0,
            max_original_price: 0,
            original_price: 0,
            max_special_price: null,
            special_price: 0,
            max_special_price: 0,
            max_saving_percentage: null,
            special_price_from: "",
            special_price_to: "",
        },
        images: [],
        videos: null,
        meta: {
            name: "",
            fk_catalog_supplier: 0,
        },
        attributes: {
            description: "",
            sku: "",
            color_position: 0,
        },
        simples: [],
        size_meta: {
            label: "",
            type: "",
            fk_catalog_attribute_set: 0,
            product_type: "simple",
            name: "",
        },
        url_key: "",
        shipment_type: "",
        shipment_info: {
            id: 0,
            shipment_type: "",
        },
        supplier_info: {
            id: 0,
            supplier_name: "",
            supplier_status: "",
        },
        visibility: false,
        weighted_availability: 0,
        created_at: "1970-01-01T00:00:00+05:30",
        activated_at: "1970-01-01T00:00:00+05:30",
        ty: {
            id: null,
            name: null,
        },
        dispatch_location_info: {
            id: null,
        },
    }
}

var newCatalogObj = function () {
    return {
        id: 0,
        sku: "",
        name: "",
        url_key: "",
        brand: {
            id: 0,
            name: "",
            url_key: "",
            is_exclusive: 0,
        },
        price_info: {
            max_price: 0,
            price: 0,
            max_original_price: 0,
            original_price: 0,
            max_special_price: null,
            special_price: 0,
            max_special_price: 0,
            max_saving_percentage: null,
            special_price_from: "",
            special_price_to: "",
        },
        images: [{}],
        meta: {
            name: "",
        },
        simples: [],
        product_group: [],
    }
}

module.exports = {
    transformData: function (data) {
        try {
            var td = JSON.parse(data);
            var boutique = {}
            boutique.data = []
            for (index in td.data) {
                sd = td.data[index].data
                boutique.data[index] = newBoutiqueObj()
                bd = boutique.data[index]

                bd.id = sd.id
                bd.sku = sd.sku
                bd.name = sd.name
                bd.status = sd.status
                bd.url_key = sd.urlKey
                bd.created_at = sd.createdAt
                bd.visibility = td.data[index].visible
                if (sd.ty != null) {
                    bd.catalog_ty = sd.ty.name
                }

                bd.attribute_set_name = {
                    id: sd.attributeSet.seqId,
                    name: sd.attributeSet.name,
                    label: sd.attributeSet.label,
                    label_en: sd.attributeSet.label_en,
                    identifier: "",
                }

                bd.brand = {
                    id: sd.brand.id,
                    name: sd.brand.name,
                    url_key: sd.brand.urlKey,
                }
                for (i in sd.categories) {
                    bd.categories[i] = {
                        id: sd.categories[i].id,
                        name: sd.categories[i].name,
                        url_key: sd.categories[i].urlKey,
                        segment: sd.categories[i].segment,
                        segment_url_key: sd.categories[i].segmentUrlKey,
                    }
                }

                for (i in sd.group) {
                    if (bd.product_group == null) {
                        bd.product_group = []
                    }
                    bd.product_group[i] = {
                        sku: sd.group[i].sku,
                        status: sd.group[i].status,
                        pet_approved: sd.group[i].petApproved,
                        color: sd.group[i].color,
                        url_key: sd.group[i].urlKey,
                    }
                }

                if (sd.rating_info != null) {
                    bd.rating_info = [sd.rating]
                }

                bd.price_info = {
                    max_price: sd.price.maxPrice,
                    price: sd.price.price,
                    max_original_price: sd.price.maxOriginalPrice,
                    original_price: sd.price.originalPrice,
                    special_price: sd.price.discountedPrice,
                    max_special_price: sd.price.discountedPrice,
                    max_saving_percentage: sd.price.maxSavingPercentage,
                    special_price_from: sd.price.specialPriceFrom,
                    special_price_to: sd.price.specialPriceTo,
                }

                for (i in sd.images) {
                    bd.images[i] = {}
                    bd.images[i].orientation = sd.images[i].orientation
                    bd.images[i].image_list = []
                    for (j in sd.images[i].imageList) {
                        bd.images[i].image_list[j] = {}
                        bd.images[i].image_list[j].image_no = sd.images[i].imageList[j].imageNo
                        bd.images[i].image_list[j].main = sd.images[i].imageList[j].main
                        bd.images[i].image_list[j].original_file_name = sd.images[i].imageList[j].originalFilename
                        bd.images[i].image_list[j].image_name = sd.images[i].imageList[j].imageName
                    }
                }

                if (sd.videos.length != 0) {
                    bd.videos = []
                    for (i in sd.videos) {
                        bd.videos[i] = {
                            id: sd.videos[i].id,
                            file_name: sd.videos[i].fileName,
                            thumbnail: sd.videos[i].thumbNail,
                        }
                    }
                }
                bd.meta.name = sd.name;
                bd.meta.fk_catalog_supplier = sd.supplier.seqId;
                for (i in sd.meta) {
                    if (sd.meta[i].attributeType != "multi_option") {
                        if (sd.meta[i].name == "dispatch_location") {
                            bd.dispatch_location_info.id = getAttrId(sd.meta[i].value)
                        }
                        if (sd.meta[i].name == "processing_time") {
                            sd.meta[i].value = parseInt(sd.meta[i].value)
                        }
                        bd.meta[sd.meta[i].name] = sd.meta[i].value
                        continue
                    }
                    if (sd.meta[i].value != null) {
                        bd.meta[sd.meta[i].name] = sd.meta[i].value.join("|")
                        continue
                    }
                    bd.meta[sd.meta[i].name] = null
                }

                for (i in sd.attributes) {
                    if (sd.attributes[i].attributeType != "multi_option") {
                        bd.attributes[sd.attributes[i].name] = sd.attributes[i].value
                        continue
                    }
                    if (sd.attributes[i].value != null) {
                        bd.attributes[sd.attributes[i].name] = sd.attributes[i].value.join("|")
                        continue
                    }
                    bd.attributes[sd.attributes[i].name] = null
                }
                bd.attributes.description = sd.description
                bd.attributes.sku = sd.sku;


                bd.size_meta = {
                    fk_catalog_attribute_set: sd.attributeSet.seqId,
                    product_type: "simple",
                }

                for (i in sd.simples) {
                    bd.simples[i] = {
                        id: sd.simples[i].id,
                        barcode_ean: sd.simples[i].barcodeEan,
                        sku: sd.simples[i].sku,
                        special_price: sd.simples[i].discountedPrice,
                        tax_percent: 0,
                        special_price_from: sd.simples[i].specialFromDate,
                        special_price_to: sd.simples[i].specialToDate,
                        sku_supplier_simple: sd.simples[i].barcodeEan,
                        max_saving_percentage: sd.simples[i].maxSavingPercentage,
                        attributes: {
                            barcode_ean: sd.simples[i].barcodeEan,
                        }
                    }
                    for (j in sd.simples[i].attribute) {
                        if (sd.simples[i].attribute[j].attributeType != "multi-option") {
                            bd.simples[i].attributes[sd.simples[i].attribute[j].name] = sd.simples[i].attribute[j].value
                            continue
                        }
                        if (sd.simples[i].attribute[j].value != null) {
                            bd.simples[i].attributes[sd.simples[i].attribute[j].name] = sd.simples[i].attribute[j].value.join("|")
                            continue
                        }
                        bd.simples[i].attributes[sd.simples[i].attribute[j].name] = null
                    }
                    bd.simples[i].meta = {
                        jabong_discount: 0,
                        price: sd.simples[i].price,
                        sku: sd.simples[i].sku,
                        ean_code: sd.simples[i].barcodeEan,
                        original_price: sd.simples[i].price,
                        seller_sku: sd.simples[i].barcodeEan,
                    }
                    for (j in sd.simples[i].meta) {
                        if (sd.simples[i].meta[j]) {
                            if (sd.simples[i].meta[j] != "multi-option") {
                                var name = sd.simples[i].meta[j].name;
                                bd.simples[i].meta[name] = sd.simples[i].meta[j].value;
                                if (name == "sh_size" || name == "apk_size" || name == "apm_size" || name == "apw_size" || name == "variation") {
                                    bd.size_meta.label = sd.simples[i].meta[j].label
                                    bd.size_meta.type = sd.simples[i].meta[j].attributeType
                                    bd.size_meta.name = sd.simples[i].meta[j].name
                                    bd.simples[i].meta[name + "_position"] = parseInt(sd.simples[i].position)
                                }
                                continue
                            }
                            if (sd.simples[i].meta[j].value) {
                                bd.simples[i].meta[sd.simples[i].meta[j].name] = sd.simples[i].meta[j].value.join("|")
                                continue
                            }
                            bd.simples[i].meta[sd.simples[i].meta[j].name] = null
                        }
                    }
                }
                if (sd.shipment.id) {
                    bd.shipment_info = {
                        id: sd.shipment.id,
                        shipment_type: sd.shipment.type,
                    }
                }
                bd.supplier_info = {
                    id: sd.supplier.seqId,
                    supplier_name: sd.supplier.slrName,
                    supplier_status: sd.supplier.status,
                }

                if (sd.ty != null) {
                    bd.ty = {
                        id: sd.ty.id,
                        name: sd.ty.name,
                    }
                }

                if (sd.sizeChart != null) {
                    bd.sizechart.id = sd.id,
                        bd.sizechart.type = sd.sizeChart.sizeChartType
                }
            }
        } catch (err) {
            console.log(err)
            console.log(err.stack);
            boutique.err = err.message
            return boutique
        }
        return boutique
    },
    transformDataSmall: function (data) {
        try {

            var td = JSON.parse(data);
            var boutique = {}
            boutique.data = []
            for (index in td.data) {
                sd = td.data[index].data
                boutique.data[index] = {}
                bd = boutique.data[index]

                bd.id = sd.id
                bd.sku = sd.sku
                bd.name = sd.name
                bd.status = sd.status

                bd.rating_info = [sd.rating]

                bd.price_info = {
                    max_price: sd.price.maxPrice,
                    price: sd.price.price,
                    max_original_price: sd.price.maxOriginalPrice,
                    original_price: sd.price.originalPrice,
                    special_price: sd.price.discountedPrice,
                    max_special_price: sd.price.discountedPrice,
                    max_saving_percentage: sd.price.maxSavingPercentage,
                    special_price_from: sd.price.specialPriceFrom,
                    special_price_to: sd.price.specialPriceTo,
                }

                bd.images = []
                for (i in sd.images) {
                    bd.images[i] = {}
                    bd.images[i].orientation = sd.images[i].orientation
                    bd.images[i].image_list = []
                    for (j in sd.images[i].imageList) {
                        bd.images[i].image_list[j] = {}
                        bd.images[i].image_list[j].image_no = sd.images[i].imageList[j].imageNo
                        bd.images[i].image_list[j].main = sd.images[i].imageList[j].main
                        bd.images[i].image_list[j].original_file_name = sd.images[i].imageList[j].originalFilename
                        bd.images[i].image_list[j].image_name = sd.images[i].imageList[j].imageName
                    }
                }

                bd.size_meta = {
                    fk_catalog_attribute_set: sd.attributeSet.seqId,
                }

                bd.simples = []
                for (i in sd.simples) {
                    bd.simples[i] = {}
                    bd.simples[i].id = sd.simples[i].id
                    bd.simples[i].sku = sd.simples[i].sku
                    bd.simples[i].special_price = sd.simples[i].discountedPrice
                    for (j in sd.simples[i].meta) {
                        bd.size_meta.label = sd.simples[i].meta[j].label
                        bd.size_meta.type = sd.simples[i].meta[j].attributeType
                        bd.size_meta.name = sd.simples[i].meta[j].name
                    }
                }

                bd.url_key = sd.urlKey
            }
        } catch (err) {
            console.log(err)
            console.log(new Error().stack)
            boutique.err = err.message
            return boutique
        }
        return boutique
    },
    transformDataCatalog: function (data) {
        try {

            var td = JSON.parse(data);
            var boutique = {}
            boutique.data = []
            for (index in td.data) {
                sd = td.data[index].data
                boutique.data[index] = newCatalogObj()
                bd = boutique.data[index]
                bd.id = sd.id
                bd.sku = sd.sku
                bd.name = sd.name
                bd.url_key = sd.urlKey
                bd.created_at = sd.createdAt

                bd.brand = {
                    id: sd.brand.id,
                    name: sd.brand.name,
                    url_key: sd.brand.urlKey,
                    is_exclusive: sd.brand.isExclusive,
                }

                bd.price_info = {
                    max_price: sd.price.maxPrice,
                    price: sd.price.price,
                    max_original_price: sd.price.maxOriginalPrice,
                    original_price: sd.price.originalPrice,
                    special_price: sd.price.discountedPrice,
                    max_special_price: sd.price.discountedPrice,
                    max_saving_percentage: sd.price.maxSavingPercentage,
                    special_price_from: sd.price.specialPriceFrom,
                    special_price_to: sd.price.specialPriceTo,
                }

                bd.images[0] = {
                    orientation: sd.image.orientation,
                    image_list: [{}],
                }
                bd.images[0].image_list[0] = {
                    image_no: sd.image.imageNo,
                    main: sd.image.main,
                    original_file_name: sd.image.originalFilename,
                    image_name: sd.image.imageName,
                }

                bd.meta.name = sd.name;
                for (i in sd.meta) {
                    if (sd.meta[i].attributeType != "multi_option") {
                        if (sd.meta[i].name == "processing_time") {
                            sd.meta[i].value = parseInt(sd.meta[i].value)
                        }
                        bd.meta[sd.meta[i].name] = sd.meta[i].value
                        continue
                    }
                    if (sd.meta[i].name == "color_family") {
                        bd.meta[sd.meta[i].name] = sd.meta[i].value[0]
                        continue
                    }
                    if (sd.meta[i].value) {
                        bd.meta[sd.meta[i].name] = sd.meta[i].value.join("|");
                        continue
                    }
                    bd.meta[sd.meta[i].name] = null;
                }

                for (i in sd.simples) {
                    bd.simples[i] = {}
                    bd.simples[i].id = sd.simples[i].id
                    bd.simples[i].sku = sd.simples[i].sku
                    bd.simples[i].size = sd.simples[i].size
                }

                for (i in sd.group) {
                    if (bd.product_group == null) {
                        bd.product_group = []
                    }
                    bd.product_group[i] = {
                        sku: sd.group[i].sku,
                        id: sd.group[i].id,
                        title: sd.group[i].name,
                        created_at: sd.group[i].createdAt,
                        url_key: sd.group[i].urlKey,
                        color: {},
                        brand_name: "",
                        price_info: {},
                        simples: [],
                        images: [],
                    }
                }

                for (i in sd.group) {
                    if (sd.group[i].color != null) {
                        bd.product_group[i].color = sd.group[i].color.join("|")
                    }
                    bd.product_group[i].color = null;
                    bd.product_group[i].brand_name = sd.group[i].brand.name;
                    bd.product_group[i].price_info = {
                        max_price: sd.group[i].priceMap.maxPrice,
                        price: sd.group[i].priceMap.price,
                        max_original_price: sd.group[i].priceMap.maxOriginalPrice,
                        original_price: sd.group[i].priceMap.originalPrice,
                        special_price: sd.group[i].priceMap.discountedPrice,
                        max_special_price: sd.group[i].priceMap.discountedPrice,
                        max_saving_percentage: sd.group[i].priceMap.maxSavingPercentage,
                        special_price_from: sd.group[i].priceMap.specialPriceFrom,
                        special_price_to: sd.group[i].priceMap.specialPriceTo,
                    }
                    for (j in sd.group[i].simples) {
                        bd.product_group[i].simples[j] = {
                            id: sd.group[i].simples[j].id,
                            sku: sd.group[i].simples[j].sku,
                            size: sd.group[i].simples[j].size,
                            special_price: sd.group[i].priceMap.discountedPrice,
                            quantity: 0,
                        }
                    }
                    bd.product_group[0].images[0] = {
                        orientation: sd.image.orientation,
                        image_list: [{}],
                    }
                    bd.product_group[0].images[0].image_list[0] = {
                        image_no: sd.image.imageNo,
                        main: sd.image.main,
                        original_file_name: sd.image.originalFilename,
                        image_name: sd.image.imageName,
                    }
                }
            }
        } catch (err) {
            console.log(err)
            console.log(new Error().stack)
            boutique.err = err.message
            return boutique
        }
        return boutique
    },
}
