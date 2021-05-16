var baseURL = 'http://127.0.0.1:8089/'

function login(data) {
    return new Promise(function(resolve, reject) {
        $.ajax({
            type: "get",
            url: baseURL + 'login',
            data:data,
            dataType: "json", //指定服务器返回的数据类型
            success: function (response) {
                if (response.code != 0) {
                    reject(response.msg)
                } else {
                    resolve(response)
                }
            },
            error:function(err){
                reject(err)
            }
        });
    })
}

function score(data) {
    return new Promise(function(resolve, reject) {
        $.ajax({
            type: "post",
            url: baseURL + 'score',
            data:data,
            dataType: "json", //指定服务器返回的数据类型
            success: function (response) {
                if (response.code != 0) {
                    reject(response.msg)
                } else {
                    resolve(response)
                }
            },
            error:function(err){
                reject(err)
            }
        })
    })
}

function rank(data) {
    return new Promise(function(resolve, reject) {
        $.ajax({
            type: "get",
            url: baseURL + 'rank',
            data:data,
            dataType: "json", //指定服务器返回的数据类型
            success: function (response) {
                if (response.code != 0) {
                    reject(response.msg)
                } else {
                    resolve(response)
                }
            },
            error:function(err){
                reject(err)
            }
        })
    })
}