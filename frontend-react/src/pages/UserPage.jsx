import { React, useEffect, useState } from "react";
import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid2";

import GameList from "../components/GameList";

import "../styles/userPage.css";

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: "#fff",
  ...theme.typography.body2,
  padding: theme.spacing(1),
  textAlign: "center",
  color: theme.palette.text.secondary,
  ...theme.applyStyles("dark", {
    backgroundColor: "#1A2027",
  }),
}));

const webServiceUrl = process.env.REACT_APP_WEB_SERVICE_URL;

function getUserData(setUserData) {
  fetch(webServiceUrl + "/getData?steamid=76561199468149920")
    .then((response) => response.json())
    .then((data) => {
      // Process the data
      console.log("Data", data);
      setUserData(data);
    })
    .catch((error) => {
      // Handle the error
    });
}

export default function UserPage() {
  const [userData, setUserData] = useState({});

  useEffect(() => {
    console.log("start");
    getUserData(setUserData);
  }, []);

  return (
    <div className="Wrapper">
      <Grid className="Panel" container spacing={6}>
        <Grid className="SidePanel" size={2}>
          <Item className="Profile">
            <img src={userData.avatar_full_url} />
            <h2 className="PersonaCard">{userData.persona_name}</h2>
          </Item>
          <Item className="StatsCard" style={{ marginTop: "30px" }}></Item>
        </Grid>
        <Grid className="MainPanel" size={10}>
          <Item className="MainCard">
            <Grid className="ChartWrapper" container spacing={2}>
              <Grid className="ChartCard" size={6}>
                <GameList userData={userData} />
              </Grid>
              <Grid className="ChartCard" size={6}>
                <div>Test</div>
              </Grid>
              <Grid className="ChartCard" size={6}>
                <div>Test</div>
              </Grid>
              <Grid className="ChartCard" size={6}>
                <div>Test</div>
              </Grid>
            </Grid>
          </Item>
        </Grid>
      </Grid>
    </div>
  );
}
