# etcd v0.4.6

./etcd -peer-addr 127.0.0.1:7001 -addr 127.0.0.1:4001 -data-dir machines/machine1 -name machine1
./etcd -peer-addr 127.0.0.1:7002 -addr 127.0.0.1:4002 -peers 127.0.0.1:7001,127.0.0.1:7003 -data-dir machines/machine2 -name machine2
./etcd -peer-addr 127.0.0.1:7003 -addr 127.0.0.1:4003 -peers 127.0.0.1:7001,127.0.0.1:7002 -data-dir machines/machine3 -name machine3

curl -L http://127.0.0.1:4001/v2/machines
curl -L http://127.0.0.1:4001/v2/keys/_etcd/machines
curl -L http://127.0.0.1:4001/v2/leader

curl -L http://127.0.0.1:4001/v2/keys/foo -XPUT -d value=bar
curl -L http://127.0.0.1:4002/v2/keys/foo

