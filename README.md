php_session_decoder
===================

PHP session encoder/decoder written in Go  

Installation
------------

Install:

    go get github.com/yvasiyarov/php_session_decoder

Getting started
---------------

Example: load php session data from redis:

    if sessionId, err := req.Cookie("frontend"); err == nil {
        if sessionData, err := redis.Get("PHPREDIS_SESSION:" + sessionId.Value); err == nil {
            decoder := php_session_decoder.NewPhpDecoder(sessionData.String())
            if sessionDataDecoded, err := decoder.Decode(); err == nil {
                //Do something with session data  
            }
        } else {
            //Can not load session - it can be expired
        }
    }

Example: Encode php session data:

    data := make(PhpSessionData)
    data["make some"] = " changes"
    encoder := NewPhpEncoder(data)
    if result, err := encoder.Encode(); err == nil {
        //Write data to redis/memcached/file/etc
    }

Some conversion details:  

1. Int, string, bool convert one to one without any surprises.  
2. Float values now has less precision then in PHP, so it can be truncated. Probably its bug, should be investigated.   
3. Array keys is always decoded as strings. Even if they were numbers in PHP. They converted back to numbers when session is encoded.  
4. PHP object instances decoded as intances of PhpObject. Please see common.go
