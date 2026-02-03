function Ajax(url, method, data, onOK, onCancel) {
    let xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState === xmlHttp.HEADERS_RECEIVED) {
            // Get the raw header string
            const headers = xmlHttp.getAllResponseHeaders();

            // Convert the header string into an array
            // of individual headers
            const arr = headers.trim().split(/[\r\n]+/);

            console.log(arr)
        }
        if (xmlHttp.readyState === 4) {
            if (xmlHttp.status >= 200 && xmlHttp.status < 300) {
                onOK(xmlHttp)
            } else {
                onCancel(xmlHttp)
            }
        }
    }
    xmlHttp.open(method, url, false);
    xmlHttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    if (App.UserToken != null) {
        console.log("send Authorization message")
        xmlHttp.setRequestHeader("Authorization", "Bearer " + App.UserToken);
    }
    xmlHttp.send(data);
}
