import React from "react";
import styled from "styled-components";

import NewRideForm from "./components/NewRideForm";
import RideList from "./components/RideList";

import "./App.css";
import NewLocationForm from "./components/NewLocationForm";
import { BrowserRouter } from "react-router-dom";

export default function App() {
  return (
    <BrowserRouter>
      <StyledContainer>
        <h1>CountryRoads</h1>
        <NewRideForm/>
        <NewLocationForm/>
        <RideList/>
      </StyledContainer>
    </BrowserRouter>

  );
}

const StyledContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;

  height: 100%;
  width: 100%;
`;
