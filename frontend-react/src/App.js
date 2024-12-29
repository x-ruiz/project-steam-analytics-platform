import logo from "./logo.svg";
import "./App.css";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import UserPage from "./pages/UserPage";
import HomePage from "./pages/HomePage";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/user/:steamid" element={<UserPage />} />
        </Routes>
      </BrowserRouter>
      {/* <div className="App"> */}
      {/* <HomePage /> */}
      {/* </div> */}
    </ThemeProvider>
  );
}

export default App;
