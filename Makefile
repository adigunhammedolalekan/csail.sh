build:
	./build.sh

compile:
	export GOOS=linux && go build -o hostgolang cmd/cmd.go

kube-deploy:
	# kubectl apply -f k8s/secrets.yml
	kubectl apply -f k8s/configs.yml
	kubectl apply -f k8s/services.yml
	kubectl apply -f k8s/statefulsets.yml
	kubectl apply -f k8s/deployments.yml
deploy:
	make build
	make k8s