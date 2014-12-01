zookeeper
=========

### Setup server

    ZK_BASE=/sgn/app/zookeeper
    mkdir -p $ZK_BASE
    cd $ZK_BASE
    wget -q http://apache.mirrors.pair.com/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz
    tar -xzf zookeeper-3.4.6.tar.gz
    
    MYID=$1
    mkdir -p $ZK_BASE/var/zookeeper/{data,conf}
    echo -n $MYID > $ZK_BASE/var/zookeeper/data/myid
    cat > $ZK_BASE/var/zookeeper/conf/zoo.cfg <<EOF
    tickTime=2000 # in ms, s-s/c-s heartbeat interval
    initLimit=10
    syncLimit=5
    dataDir=$ZK_BASE/var/zookeeper/data
    dataLogDir=$ZK_BASE/var/zookeeper/dataLog
    clientPort=2181
    server.1=192.168.12.11:2888:3888
    server.2=192.168.12.12:2888:3888
    server.3=192.168.12.13:2888:3888
    EOF
    
    $ZK_BASE/zookeeper-3.4.6/bin/zkServer.sh start-foreground $ZK_BASE/var/zookeeper/conf/zoo.cfg

### CLI

    go get github.com/mmcgrana/zk
