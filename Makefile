tf-fmt:
	cd terraform && terraform fmt

tf-init:
	cd terraform && terraform init

tf-plan:
	cd terraform && terraform plan

tf-apply:
	cd terraform && terraform apply

start-server:
	cd web-service-go && go run .

start-frontend:
	cd frontend-react && npm start