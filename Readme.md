
CURL CHAT
=========

Bidirectional chat over curl

## Usage

```shell
curlchat -h
  -baseurl string
        the base url of the service
         (default "http://localhost:8080/")
  -httpaddr string
        the address/port to listen on for http
        use :<port> to listen on all addresses
         (default "localhost:8080")
```

## Building

In order to build the project, just use:
```shell
go build
```

## Server Install

See `scripts/install.sh`

## Server Deployment

See `scripts/deploy.sh`

## Todo

* html formatting
* initial history
