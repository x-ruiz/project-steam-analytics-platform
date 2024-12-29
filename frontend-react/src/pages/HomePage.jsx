import { React, useState } from "react";
import { useNavigate } from "react-router-dom";

import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import ArrowForwardIos from "@mui/icons-material/ArrowForwardIos";

import "../styles/HomePage.css";

export default function HomePage() {
  const [steamId, setSteamId] = useState("");

  const navigate = useNavigate();
  const handleClick = (steamId) => {
    navigate(`/user/${steamId}`);
  };

  return (
    <div className="Wrapper">
      <div className="Center">
        <Box
          component="form"
          sx={{ "& > :not(style)": { m: 1, width: "25ch" } }}
          noValidate
          autoComplete="off"
        >
          <TextField
            id="outlined-basic"
            label="SteamID"
            variant="outlined"
            onChange={(e) => setSteamId(e.target.value)}
          />
        </Box>
        <Button
          variant="contained"
          endIcon={<ArrowForwardIos />}
          onClick={() => handleClick(steamId)}
        >
          View Stats
        </Button>
      </div>
    </div>
  );
}
