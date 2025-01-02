resource "google_bigquery_dataset" "dataset" {
  dataset_id  = var.dataset_name
  description = var.dataset_description
  location    = "us-central1"
}

resource "google_bigquery_table" "user_table" {
  dataset_id = "main"
  table_id   = "t_user_table"

  description         = "Table to store user data"
  deletion_protection = false
  schema = jsonencode([
    {
      "name" : "timestamp",
      "type" : "TIMESTAMP",
      "mode" : "REQUIRED",
      "description" : "Timestamp of data insert"
    },
    {
      "name" : "steam_id",
      "type" : "STRING",
      "mode" : "REQUIRED",
      "description" : "Steamid of the user"
    },
    {
      "name" : "persona_name",
      "type" : "STRING",
      "mode" : "REQUIRED",
      "description" : "Persona Name of the user"
    },
    {
      "name" : "game_count",
      "type" : "INTEGER",
      "mode" : "REQUIRED",
      "description" : "Total number of games owned by the user account"
    },
    {
      "name" : "games",
      "type" : "RECORD",
      "mode" : "REPEATED",
      "description" : "List of games and metadata of games",
      "fields" : [
        {
          "name" : "appid",
          "type" : "INTEGER",
          "mode" : "REQUIRED"
        },
        {
          "name" : "name",
          "type" : "STRING",
          "mode" : "REQUIRED"
        },
        {
          "name" : "playtime_forever",
          "type" : "INTEGER",
          "mode" : "REQUIRED"
        },
        {
          "name" : "img_icon_url",
          "type" : "STRING",
          "mode" : "NULLABLE"
        },
        {
          "name" : "img_logo_url",
          "type" : "STRING",
          "mode" : "NULLABLE"
        }
      ]
    }
  ])
}

resource "google_bigquery_table" "steam_lifetime_playtime_view" {
  table_id            = "v_steam_lifetime_playtime"
  dataset_id          = "main"
  deletion_protection = false
  view {
    query          = file("${path.module}/sql/steam_lifetime_playtime.sql")
    use_legacy_sql = false
  }
}
