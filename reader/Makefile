build:
	GOOS=linux go build -o app
	docker build -t docker-registry-default.apps.dev-cefcfco.com/rma-7x24/websocket-client-mock .
	rm -f app