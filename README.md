# Dasher

Dasher is a json rest api written in Golang.

This app is just a playground to learn how to better write idiomatic Go.

Trying to stick to the standard library as much as possible.  Added negroni to help with middleware handlers (logger).  Playing with prometheus client for adding instrumentation to handlers.


Requires v1.8 since using new http.Shutdown()



## Configuration

```
{
    "version" : "1.0.0",
    "address": "0.0.0.0:3000",
    "environment": "dev",
    "shutdownTimeout": 5,
    "debug": false,
    "analyticsOAuth": {
        "pemFilePath": "./sample.pem",
        "serviceEmail": "readonly@sample.iam.gserviceaccount.com",
        "scopes": [ 
            "https://www.googleapis.com/auth/analytics.readonly" 
        ]
    }
}
```

### TODO
- Add tests
- Looking for ideas to include gRPC
- Build additional middleware components (security checks, login)
- Learn how to use 1.8 Context pkg
- Build a slick frontend
- Get oauth login working 
- Test go/dep for /vendor
- Learn more!!



