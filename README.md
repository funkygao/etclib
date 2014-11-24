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
           
