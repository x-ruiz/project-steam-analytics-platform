import { React, useEffect, useState } from "react";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import ArrowForwardIos from "@mui/icons-material/ArrowForwardIos";

import "./styles/SteamId.css";

const webServiceUrl = process.env.REACT_APP_WEB_SERVICE_URL;

function getUserData(setUsername, username) {
  fetch(webServiceUrl + "/getSteamId?username=" + username)
    .then((response) => response.json())
    .then((data) => {
      // Process the data
      setUsername(data["response"]["steamid"]);
    })
    .catch((error) => {
      // Handle the error
    });
}

export default function SteamId() {
  const [username, setUsername] = useState("");
  const [usernameInput, setUsernameInput] = useState("");

  const handleGetSteamId = (usernameInput) => {
    // http get steam id
    getUserData(setUsername, usernameInput);
  };

  return (
    <div className="Center">
      <Box
        component="form"
        sx={{ "& > :not(style)": { m: 1, width: "25ch" } }}
        noValidate
        autoComplete="off"
      >
        <TextField
          id="outlined-basic"
          label="Username"
          variant="outlined"
          onChange={(e) => setUsernameInput(e.target.value)}
        />
      </Box>
      <Button
        variant="contained"
        endIcon={<ArrowForwardIos />}
        onClick={() => handleGetSteamId(usernameInput)}
      >
        Get SteamID
      </Button>
      <div>{username}</div>
    </div>
  );
}
