
# 🦅 Falcon : A dynamic reverse-proxy

A dynamic reverse-proxy server built using Go which can be used as Load Balancer and API Gateway.


[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://choosealicense.com/licenses/mit/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yashkundu/falcon.svg)](https://pkg.go.dev/github.com/yashkundu/falcon)

## Contents

- [Installation](#installation)
- [Features](#features)
- [Configuration](#configuration)
- [API referencr](#api-reference)
- [License](#license)

## Installation

Falcon can be installed directly via go install. To install the latest version:

```bash
  go install github.com/yashkundu/falcon@latest
```
    
To install a specific release:
```bash
  go install github.com/yashkundu/falcon@v1.0.0
```

## Running Falcon

To start the proxy-server run :

```bash
  falcon -config=<configFilePath>
```


## Features

- L7 Load balancing
- Rate Limiting
- Dynamic Backend URLs
- Can be used as an API Gateway

## Configuration

The config file is in bin/config/config.toml\
Refer to [TOML Specificataion](https://toml.io/en/).

```toml
# the main configs of the reverse-proxy
[core]

# the port where reverse-proxy will listen
# default 80
listen=8000

# the port where the server api will listen
# refer to the below api reference section
# default 9900
apiport=9900

# Enable server stats route
# default - false
enableServerStats=false

# The maximum connections that the reverse-proxy will allow at a particular time
# 0 is infinite
limitMaxConn=0

# if it is 0, there is no timeout, in seconds
readTimeout=0

# if it is 0, there is no timeout, in seconds
writeTimeout=0

# if it is 0, there is no timeout, in seconds
idleTimeout=0

# enable this to implement rate limiting
[limitReq]
enable=false
# in millisecond
interval=1000
frequency=100

# config of proxy
[proxy]

[[proxy.routes]]
endpoint="/hello"
# match (specifies how routes are mathed to the incoming request)
# match - [0 - exact, 1 - prefix, 2 - regex ]
# default - exact
match=1
# balancer (specifies which load balancing algo to be used in case of multiple backends)
# balancer - [0 - roundrobin, 1 - random, 2 - weighted-roundrobin ]
balancer=0

[[proxy.routes.backends]]
# url should be in proper format <schema>://<host>:<port>
url="http://localhost:3000"


# Add varName to enable dynamically changing urls of the particular backend
# varNames for all the backends should be unique othewise only the last backend 
# specified with a particular varName will be dynamic
varName="x1"


[[proxy.routes]]
endpoint="/world"
match=1

[[proxy.routes.backends]]
url="http://localhost:3005"


```
## API Reference

#### Get current requests handled at this moment

```bash
  GET <proxyHostName>:9900/apiStatus/reqCount
```

#### Dynamically update the backend url

```bash
  POST <proxyHostName>:9900/apiStatus/backendChange
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `varName` | `string` | varName of the backend            |
| `url`      | `string` | the new backendUrl                |




## License

[MIT](https://choosealicense.com/licenses/mit/)
