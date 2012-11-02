function setCurrent(url) { 
        var host = "null";

        url = window.location.href;
        var regex = /.*\:\/\/[^\/]*\/admin\/([^\/\-\?\#]*).*/;
        var match = url.match(regex);
        if(typeof match != "undefined"
                        && null != match)
                host = match[1];
        return host;
}