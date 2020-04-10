build:
	./build.sh

compile:
	export GOOS=linux && go build -o hostgolang cmd/cmd.go

test:
	go generate ./...
	go test ./...

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

kube-local:
	kubectl apply -f k8s/local/storage-class.yml
	kubectl apply -f k8s/local/secrets.yml
	kubectl apply -f k8s/local/configs.yml
	kubectl apply -f k8s/local/services.yml
	kubectl apply -f k8s/local/statefulsets.yml
	kubectl apply -f k8s/local/deployments.yml
	kubectl apply -f k8s/local/tls.yml
	kubectl apply -f k8s/local/ingress.yml

ingress:
	# install ingress
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.1/deploy/static/provider/cloud-generic.yaml

delete-ingress:
	kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
	kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.26.1/deploy/static/provider/cloud-generic.yaml

cert-manager:
	# install cert-manager
	kubectl create namespace cert-manager
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.0/cert-manager.yaml