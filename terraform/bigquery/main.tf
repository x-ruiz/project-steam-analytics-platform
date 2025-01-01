resource "google_bigquery_dataset" "dataset" {
  dataset_id  = var.dataset_name
  description = var.dataset_description
  location    = "us-central1"
}