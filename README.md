etclib
======

shared config service lib integrated with etcd

### services

* faed
* actord
* maintain mode of global | kingdom_N

### health check

### nodes

        /
        ├── proj
        │   ├── dw
        |   |   ├── maintain
        |   |   |   ├── global
        |   |   |   ├── kingdom_1
        |   |   |   └── kingdom_N
        |   |   |   
        |   |   └── node
        |   |       ├── act
        |   |       |   ├── 12.10.0.1:9001
        |   |       |   ├── 12.10.0.2:9002
        |   |       |   └── 12.10.0.5:9001
        |   |       |   
        |   |       └── fae
        |   |           ├── 12.10.0.1:9898
        |   |           ├── 12.10.0.2:9898
        |   |           └── 12.10.0.3:9898
        |   |   
        │   └── xxx
           


### watch fails

    GET /v2/keys/proj/dw/node/act?consistent=true&recursive=true&wait=true HTTP/1.1 
    
    HTTP/1.1 200 OK
    Content-Type: application/json
    X-Etcd-Cluster-Id: 7e27652122e8b2ae
    X-Etcd-Index: 2506
    X-Raft-Index: 348157
    X-Raft-Term: 1
    Date: Wed, 26 Nov 2014 10:25:34 GMT
    Transfer-Encoding: chunked
    
    [11/26/14 18:25:34 CST] [EROR] (    service.go:175) watch node[act]: unexpected end of JSON input

### todo

    https://github.com/kelseyhightower/confd
    fae node has weight
