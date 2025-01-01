module "bigquery" {
  source = "./bigquery"
  
  dataset_name = "main"
  dataset_description = "Main dataset to hold user tables"
}