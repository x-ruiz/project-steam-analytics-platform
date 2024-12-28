# TO RUN

cd web-service-gin && go run .

# Deploying

https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service

gcloud projects add-iam-policy-binding steam-analytics-platform \
 --member=serviceAccount:789838811617-compute@developer.gserviceaccount.com \
 --role=roles/cloudbuild.builds.builder

### To setup Github Actions + Cloud Run

https://cloud.google.com/blog/products/devops-sre/deploy-to-cloud-run-with-github-actions/

# STEAM API Methods Used

GET | https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=XXXXXXXXXXXXXXXXX&vanityurl=Kodiris
GET | https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=XXXXXXXXXXXXXXXXX&steamid=76561197960434622&format=json
