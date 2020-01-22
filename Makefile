build:
	./build.sh

compile:
	export GOOS=linux && go build -o hostgolang cmd/cmd.go

kube-deploy:
	kubectl apply -f k8s/secrets.yml
	kubectl apply -f k8s/configs.yml
	kubectl apply -f k8s/services.yml
	kubectl apply -f k8s/tls.yml
	kubectl apply -f k8s/statefulsets.yml
	kubectl apply -f k8s/deployments.yml
deploy:
	make build
	make k8s

ingress-setup:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.1/deploy/static/mandatory.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.1/deploy/static/provider/cloud-generic.yaml

cert-manager:
	kubectl apply -f k8s/certs-test.yml