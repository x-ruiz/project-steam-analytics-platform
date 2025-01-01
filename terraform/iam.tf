resource "google_service_account" "bigquery_service_account" {
  account_id = "bigquery-service-account"
  display_name = "Bigquery Service Account"
}