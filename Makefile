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

start-server:
	cd web-service-go && go run .

start-frontend:
	cd frontend-react && npm start