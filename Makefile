build:
	./build.sh
deploy:
	make build
	kubectl apply -f k8s/configs.yml
	kubectl apply -f k8s/services.yml
	kubectl apply -f k8s/statefulsets.yml
	kubectl apply -f k8s/deployments.yml