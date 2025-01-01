terraform {
  backend "gcs" {
    bucket      = "steam-analytics-platform-tf-state"
    prefix      = "tf-state"
    credentials = "steam-analytics-platform-1440e52a1785.json"
  }
}