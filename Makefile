kind-cluster:
	kind create cluster --name=kxc --config=kind-config.yaml

k8s-deploy:
	kubectl apply -f k8s/ingress/

run-kxc-api:
	(cd kxc-api && go run cmd/main.go)

run-kxc-operator:
	(cd kxc-operator && make run ENABLE_WEBHOOKS=false)