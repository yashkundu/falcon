
# 🦅 Falcon : A dynamic reverse-proxy

A dynamic reverse-proxy server built using Go which can be used as Load Balancer and API Gateway.


[![MIT License](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://choosealicense.com/licenses/mit/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yashkundu/falcon.svg)](https://pkg.go.dev/github.com/yashkundu/falcon)

## Contents

- [Features](#features)
- [Quick start](#quick-start)
- [Benchmarks](#benchmarks)
- [Configuration examples](#configuration-examples)
    - [Reverse proxy with custom headers](#reverse-proxy-with-custom-headers)
    - [API gateway with version-path mapping](#api-gateway-with-version-path-mapping)
    - [Secured load balancer with weighted backends](#secured-load-balancer-with-weighted-backends)
    - [TCP proxy as a ssh relay server](#tcp-proxy-as-a-ssh-relay-server)
    - [Scooter API endpoint](#scooter-api-endpoint)
    - [Scooter prometheus endpoint](#scooter-prometheus-endpoint)
- [Migrate from nginx to scooter](#migrate-from-nginx-to-scooter)
## Installation

Falcon can be installed directly via go install. To install the latest version:

```bash
  go install github.com/yashkundu/falcon@latest
```
    
To install a specific release:
```bash
  go install github.com/yashkundu/falcon@v0.0.1
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Features

- L7 Load balancing
- Rate Limiting
- Dynamic Backend URLs
- Can be used as an API Gateway

<p align="right">(<a href="#readme-top">back to top</a>)</p>
## License

[MIT](https://choosealicense.com/licenses/mit/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>