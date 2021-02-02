
NAME                  = es-kind
ELASTICSEARCH_VERSION = 7.10.2
TIMEOUT               = 1200s

.PHONY: kind/start
kind/start:
	kind create cluster --name $(NAME)

.PHONY: kind/stop
kind/stop:
	kind delete cluster --name $(NAME)


.PHONY: helm/repo/add
helm/repo/add:
	helm repo add elastic https://helm.elastic.co

.PHONY: helm/elasticsearch/install
helm/elasticsearch/install:
	kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
	helm install elasticsearch --version $(ELASTICSEARCH_VERSION) --wait --timeout=$(TIMEOUT) -f ./charts/elasticsearch/values.yaml elastic/elasticsearch

.PHONY: helm/elasticsearch/uninstall
helm/elasticsearch/uninstall:
	helm delete elasticsearch


.PHONY: docker/start
docker/start:
	docker-compose -f ./deployments/docker-compose.yml up -d

.PHONY: docker/down
docker/down:
	docker-compose -f ./deployments/docker-compose.yml down

.PHONY: docker/stop
docker/stio:
	docker-compose -f ./deployments/docker-compose.yml stop
