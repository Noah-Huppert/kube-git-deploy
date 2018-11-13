.PHONY: api etcd etcdctl

ETCD_DATA_DIR=${PWD}/container-data/etcd

# run api
api:
	go run main.go

# run etcd
etcd:
	mkdir -p "${ETCD_DATA_DIR}"
	docker run \
		--net host \
		--rm \
		-it \
		-p 2379:2379 \
		-p 2380:2380 \
		--volume=${ETCD_DATA_DIR}:/etcd-data \
		quay.io/coreos/etcd:latest \
			/usr/local/bin/etcd \
			--data-dir=/etcd-data --name node1 \
			--initial-advertise-peer-urls http://localhost:2380 \
			--listen-peer-urls http://localhost:2380 \
			--advertise-client-urls http://localhost:2379 \
			--listen-client-urls http://localhost:2379 \
			--initial-cluster node1=http://localhost:2380

# run etcdctl
etcdctl:
	docker run \
		--net host \
		--rm \
		-it \
		quay.io/coreos/etcd:latest \
			/usr/local/bin/etcdctl \
			--endpoints http://localhost:2380 \
			${ARGS}
