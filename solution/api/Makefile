e2e-up:
	docker-compose -f ../docker-compose.e2e.yaml up -d

e2e-down:
	docker-compose -f ../docker-compose.e2e.yaml down --volumes

e2e-test:
	docker-compose -f ../docker-compose.e2e.yaml up -d --build
	go test -v ./tests/e2e
	docker-compose -f ../docker-compose.e2e.yaml down --volumes