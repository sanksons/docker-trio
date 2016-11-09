/*
 * HttpService.js - Promise based http service.
 *
 * Use this throughout the project.
 * Can be used statically by typing HttpService.serviceName(params)
 *
 *
 * Returns for all modules is a dictionary object consisting:
 * statusCode, data, location, headers, rawHeaders and raw (raw response)
 *
 * NOTE: Please do not use require() for this. Its available project wide
 * with the static variable HttpService
 */

var http = require('http');
var Promise = require('promise');
var UUID = require('node-uuid');
var zlib = require('zlib');

function HttpServiceException(message) {
  this.message = message;
  this.name = "HttpServiceException";
  this.toString = function () {
    return this.name + ": " + this.message
  }
}

function generateIds() {
  return {
    RequestID: UUID.v4(),
    TransactionID: UUID.v1(),
  };
}

function addContentAttr(options, postData, method, reqHeaders, contentType, requestIds) {
  if (requestIds == undefined) {
    requestIds = generateIds();
  }
  if (typeof (options) == "string") {
    var url = options.replace("http://", "")
    var basename = url.split("/")[0]
    var path = url.replace(basename, "")
    var hostname = basename
    var port = 80
    var hostArray = basename.split(":")
    if (contentType == undefined) {
      contentType = "application/json"
    }
    if (hostArray.length == 2) {
      hostname = hostArray[0]
      port = parseInt(hostArray[1])
    }
    if (contentType == "multipart/form-data") {
      opt = {
        method: method,
        hostname: hostname,
        port: port,
        path: path,
        headers: postData.getHeaders()
      }
      return opt
    }

    opt = {
      method: method,
      hostname: hostname,
      port: port,
      path: path,
      headers: {
        'Content-Type': contentType,
        'Content-Length': postData.length,
        'RequestID': requestIds.RequestID,
        'TransactionID': requestIds.TransactionID,
        'RequestSource':'StylokoMapper',
      }
    }
    if (typeof reqHeaders !== 'undefined') {
      for (var i in reqHeaders) {
        opt['headers'][i] = reqHeaders[i]
      }
    }
    return opt
  } else {
    if (options["header"] != undefined) {
      options["header"] = {
        'Content-Type': contentType,
        'Content-Length': postData.length,
        'RequestID': requestIds.RequestID,
        'TransactionID': requestIds.TransactionID,
        'RequestSource':'StylokoMapper',
      }
    }
  }
  return opt
}

var genHeaders = function (options, requestIds, reqHeaders) {
  if (requestIds == undefined) {
    requestIds = generateIds();
  }
  var headers = {
    'RequestID': requestIds.RequestID,
    'TransactionID': requestIds.TransactionID,
    'Content-Type': "application/json",
    'RequestSource':'StylokoMapper',
  }
  if (typeof (options) == "string") {
    var url = options.replace("http://", "")
    var basename = url.split("/")[0]
    var path = url.replace(basename, "")
    var hostname = basename
    var port = 80
    var hostArray = basename.split(":")
    if (hostArray.length == 2) {
      hostname = hostArray[0]
      port = parseInt(hostArray[1])
    }
    opt = {
      method: "GET",
      hostname: hostname,
      port: port,
      path: path,
      headers: headers,
    }
    if (reqHeaders!=undefined){
        for (var i in reqHeaders) {
          opt.headers[i] = reqHeaders[i]
        }
    }
    opt.headers["Accept-Encoding"]="gzip";
    return opt
  } else {
      options['headers'] = headers;
      if (reqHeaders!=undefined){
          for (var i in reqHeaders) {
            options.headers[i] = reqHeaders[i]
          }
      }
      options.headers["Accept-Encoding"]="gzip";
      return options
  }
}

module.exports = {
  get: function (options, reqHeaders, requestIds) {
    if (!options) {
      throw new HttpServiceException("No options or URL provided");
    } else {
      options = genHeaders(options, requestIds, reqHeaders)
      var promise = new Promise(function (resolve, reject) {
        http.get(options, (res) => {
          var body = '';
          if (res.headers['content-encoding']=="gzip"){
              gz = zlib.createGunzip();
              res.pipe(gz);
              output = gz;
          }else{
              output = res;
          }
          output.on('data', (chunk) => {
            body += chunk;
          });
          output.on('end', () => {
            response = {
              statusCode: res.statusCode,
              data: body,
              location: res.location,
              headers: res.headers,
              rawHeaders: res.rawHeaders,
              raw: res,
            }
            resolve(response)
          })
        }).on('error', (err) => {
          reject(err)
        })
      })
      return promise
    }
  },
  post: function (options, postData, flag, requestIds) {
    if (!options) {
      throw new HttpServiceException("No options or URL provided");
    } else {
      if (typeof postData['headers'] !== 'undefined') {
        var extendHeaders = postData['headers'];
        delete postData['headers'];
      }
      if (flag) {
        postData = postData.data;
      }

      if (typeof (postData) != "string") {
        postData = JSON.stringify(postData)
      }
      options = addContentAttr(options, postData, "POST", extendHeaders, undefined, requestIds);
      var promise = new Promise(function (resolve, reject) {
        var req = http.request(options, (res) => {
          var body = ""
          res.setEncoding('utf8');
          res.on('data', (chunk) => {
            body += chunk;
          });
          res.on('end', () => {
            response = {
              statusCode: res.statusCode,
              data: body,
              location: res.location,
              headers: res.headers,
              rawHeaders: res.rawHeaders,
              raw: res,
            }
            resolve(response)
          })
        }).on('error', (err) => {
          reject(err)
        });
        req.write(postData);
        req.end();
      })
    }
    return promise
  },

  put: function (options, putData, flag, requestIds) {
    if (!options) {
      throw new HttpServiceException("No options or URL provided");
    } else {
      if (typeof putData['headers'] !== 'undefined') {
        var extendHeaders = putData['headers'];
        delete putData['headers'];
      }
      if (flag) {
        putData = putData.data;
      }
      if (typeof (putData) != "string") {
        putData = JSON.stringify(putData)
      }
      options = addContentAttr(options, putData, "PUT", extendHeaders, undefined, requestIds);
      var promise = new Promise(function (resolve, reject) {
        var req = http.request(options, (res) => {
          var body = ""
          res.setEncoding('utf8');
          res.on('data', (chunk) => {
            body += chunk;
          });
          res.on('end', () => {
            response = {
              statusCode: res.statusCode,
              data: body,
              location: res.location,
              headers: res.headers,
              rawHeaders: res.rawHeaders,
              raw: res
            };
            resolve(response)
          })
        }).on('error', (err) => {
          reject(err)
        });
        req.write(putData);
        req.end();
      });
    }
    return promise
  },
  postFile: function (options, postData, requestIds) {
    if (!options) {
      throw new HttpServiceException("No options or URL provided");
    } else {
      options = addContentAttr(options, postData, "post", undefined, "multipart/form-data", requestIds);
      var promise = new Promise(function (resolve, reject) {
        var req = http.request(options, (res) => {
          var body = "";
          res.setEncoding('utf8');
          res.on('data', (chunk) => {
            body += chunk;
          });
          res.on('end', () => {
            response = {
              statusCode: res.statusCode,
              data: body,
              location: res.location,
              headers: res.headers,
              rawHeaders: res.rawHeaders,
              raw: res
            };
            resolve(response)
          })
        }).on('error', (err) => {
          reject(err)
        });
        postData.pipe(req);
      })
    }
    return promise
  },
  delete: function (options, requestIds) {
    if (!options) {
      throw new HttpServiceException("No options or URL provided");
    } else {
      var data = "";
      options = addContentAttr(options, data, "DELETE", undefined,  undefined, requestIds);
      var promise = new Promise(function (resolve, reject) {
        var req = http.request(options, (res) => {
          var body = ""
          res.setEncoding('utf8');
          res.on('data', (chunk) => {
            body += chunk;
          });
          res.on('end', () => {
            response = {
              statusCode: res.statusCode,
              data: body,
              location: res.location,
              headers: res.headers,
              rawHeaders: res.rawHeaders,
              raw: res
            };
            resolve(response)
          })
        }).on('error', (err) => {
          reject(err)
        });
        req.end();
      });
    }
    return promise
  },
  genRequestIds: function () {
    return generateIds();
  },
};
