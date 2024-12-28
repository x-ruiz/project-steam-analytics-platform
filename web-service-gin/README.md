# TO RUN

cd web-service-gin && go run .

# Deploying

https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service

gcloud projects add-iam-policy-binding steam-analytics-platform \
 --member=serviceAccount:789838811617-compute@developer.gserviceaccount.com \
 --role=roles/cloudbuild.builds.builder
