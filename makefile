# Set the name
NAME := webhook-proxy

# Set the shell
SHELL := /bin/bash

.PHONY: build
build:
	docker build -t webhookproxy -f build/Dockerfile .

.PHONY: savetar
savetar:
	docker save webhookproxy > webhookproxy.tar

.PHONY: exporttar
exporttar: webhookproxy.tar
	microk8s.ctr image import webhookproxy.tar

.PHONY: removetar
removetar: webhookproxy.tar
	rm -rf webhookproxy.tar

.PHONY: deploy2k8s
deploy2k8s:
	kustomize build ./deploy | kubectl apply -f -

.PHONY: undeploy2k8s
undeploy2k8s:
	kustomize build ./deploy | kubectl delete -f -

.PHONY: build2deploy
build2deploy: build savetar exporttar deploy2k8s
