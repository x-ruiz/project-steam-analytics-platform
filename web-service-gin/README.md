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

# Endpoints

http://localhost:8080/getPlaytime?steamid=76561198305662842
https://project-steam-analytics-platform-789838811617.us-central1.run.app/getData?username=kodiris

# TODO:

1. Set up React App
2. Set up go web service
3. Create starter api endpoints
4. Create sync with bigquery endpoint
5. Integrate with frontend POC
