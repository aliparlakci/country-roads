import React from "react";
import NewRideForm from "./components/NewRideForm";
import styled from "styled-components";

import "./App.css";

const StyledContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  height: 100%;
  width: 100%;
`

const StyledHeading = styled.h1`

`

function App() {
  return (
    <StyledContainer>
      <StyledHeading>CountryRoads</StyledHeading>
      <NewRideForm />
    </StyledContainer>
  );
}

export default App;
