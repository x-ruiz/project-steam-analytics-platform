provider "google" {
  credentials = file("steam-analytics-platform-1440e52a1785.json")
  project     = "steam-analytics-platform"
  region      = "us-central1"
}