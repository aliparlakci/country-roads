import React from "react";
import styled from "styled-components";

import NewRideForm from "./components/NewRideForm";
import RideList from "./components/RideList";

import "./App.css";

export default function App() {
  return (
    <StyledContainer>
      <h1>CountryRoads</h1>
      <NewRideForm />
      <RideList />
    </StyledContainer>
  );
}

const StyledContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  height: 100%;
  width: 100%;
`;
