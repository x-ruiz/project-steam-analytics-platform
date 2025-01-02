import { React, useEffect, useState } from "react";
import { DataGrid } from "@mui/x-data-grid";
import Paper from "@mui/material/Paper";

const columns = [
  { field: "id", headerName: "ID", width: 70 },
  {
    field: "img_icon_url",
    type: "actions",
    width: 70,
    cellClassName: "actions",
    renderCell: (params) => <img src={params.value} />,
  },
  { field: "name", headerName: "Name", width: 130 },
  { field: "playtime", headerName: "Playtime", width: 130 },
  { field: "playtime_forever", headerName: "Minutes", width: 130 },
];

const paginationModel = { page: 0, pageSize: 5 };

const sortModel = [
  { field: "playtime_forever", sort: "desc" }, // sort by playtime minutes in ascending order
];

function minutesToHours(playtime) {
  const hours = Math.floor(playtime / 60);
  const minutes = playtime % 60;
  return { hours, minutes };
}

export default function GameList({ userData }) {
  const rows = [];
  for (var i in userData["games"]) {
    var game = userData["games"][i];

    // Set row id as app id
    game["id"] = game["appid"];

    // playtime processing
    var playtime = minutesToHours(game["playtime_forever"]);

    game["playtime"] = `${playtime.hours}h ${playtime.minutes}m`;
    rows.push(game);
  }

  return (
    <Paper sx={{ height: 400, width: "100%" }}>
      <DataGrid
        rows={rows}
        columns={columns}
        initialState={{ pagination: { paginationModel } }}
        pageSizeOptions={[5, 10]}
        // checkboxSelection
        sx={{ border: 0 }}
        sorting={true}
        sortModel={sortModel}
      />
    </Paper>
  );
}
