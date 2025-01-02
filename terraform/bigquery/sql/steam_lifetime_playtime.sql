WITH latest_rows AS (
  SELECT
    steam_id,
    persona_name,
    games,
    timestamp,
    ROW_NUMBER() OVER (PARTITION BY steam_id ORDER BY timestamp DESC) AS row_num
  FROM
    `steam-analytics-platform.main.t_user_table`
),
filtered_rows AS (
  SELECT
    steam_id,
    persona_name,
    games,
    timestamp
  FROM
    latest_rows
  WHERE
    row_num = 1
),
game_playtime AS (
  SELECT
    steam_id,
    persona_name,
    games.appid,
    games.name,
    games.playtime_forever AS playtime,
    SUM(games.playtime_forever) OVER (PARTITION BY steam_id) AS lifetime_playtime,  -- Fix to get the total playtime per user
    (games.playtime_forever / SUM(games.playtime_forever) OVER (PARTITION BY steam_id)) AS playtime_percentage
  FROM
    filtered_rows
  CROSS JOIN
    UNNEST(games) AS games
)
SELECT
  steam_id,
  persona_name,
  appid,
  name,
  playtime,
  lifetime_playtime,
  playtime_percentage
FROM
  game_playtime
ORDER BY
  playtime_percentage DESC;
