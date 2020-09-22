# IPAddr Microservice
This useless microservice allow you to retrieve IP addresses from your server.

## Installation

To install and make `ipaddr` works you have just to type the following commands

```bash
go get github.com/julienschmidt/httprouter
go build ipaddr.go
```

## Usage

`ipaddr` is relatively simple to use.

You can optionally specify the port where you want to bind the service to, as shown in the usage below.

```bash
Usage of ./ipaddr:
  -port int
        Bind the microservice to the specified port. (default 8080)
```

To use the IPAddress Microservice, just try to curl the following endopoints:

```
http://example-server.com:8080/ipa    → show interfaces with IPv4 and IPv6 addresses
http://example-server.com:8080/ipa/4  → show interfaces with IPv4 addresses
http://example-server.com:8080/ipa/6  → show interfaces with IPv6 addresses
```

And, that's it :)