resource "google_service_account" "bigquery_service_account" {
  account_id   = "bigquery-service-account"
  display_name = "Bigquery Service Account"
}


resource "google_project_iam_binding" "bigquery_admin_policy_binding" {
  project = local.project_id # Replace with your project ID
  role    = "roles/bigquery.admin"
  members = [
    "serviceAccount:bigquery-service-account@steam-analytics-platform.iam.gserviceaccount.com"
  ]
}
