# Terraform
tf-fmt:
	@bash -c 'shopt -s globstar; \
	for dir in terraform/**/*/; do \
		if [ -d $$dir ]; then \
			cd $$dir && terraform fmt && cd -; \
		fi \
	done'

tf-init:
	cd terraform && terraform init

tf-plan:
	cd terraform && terraform plan

tf-apply:
	cd terraform && terraform apply

# Go
start-server:
	cd go-api-service && go run .

go-get-pckgs:
	cd go-api-service && go get .

# React
start-frontend:
	cd react-client && npm start