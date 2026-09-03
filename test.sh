echo "=== Running Safe Write Test =="
docker exec -e ETCDCTL_API=3 etcd-legacy etcdctl --endpoints=http://envoy:9090 --insecure-transport=true put proxy "Working"

echo "=== Running Blocked Exploit Test ==="
docker exec -e ETCDCTL_API=3 etcd-legacy etcdctl --endpoints=http://envoy:9090 --insecure-transport=true alarm list
