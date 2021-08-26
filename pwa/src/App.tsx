import React from "react";
import styled from "styled-components";

import NewRideForm from "./components/NewRideForm";
import RideList from "./components/RideList";

import "./App.css";
import NewLocationForm from "./components/NewLocationForm";

export default function App() {
  return (
    <StyledContainer>
      <h1>CountryRoads</h1>
      <NewRideForm />
      <NewLocationForm />
      <RideList />
    </StyledContainer>
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
