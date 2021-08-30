import React from "react";
import styled from "styled-components";

import RideList from "../components/RideList";
import useQuery from "../hooks/useQuery";
import { IRideQuery } from "../hooks/useRides";
import RideFilter from "../components/RideFilter";

export default function MainView() {
  const params = useQuery();
  const query: IRideQuery = {
    type: params.get("type"),
    direction: params.get("direction"),
    destination: params.get("destination"),
    startDate: params.get("start_date"),
    endDate: params.get("end_date"),
  };

  return (
    <StyledContainer>
      <ColumnView>
        <RideFilter />
      </ColumnView>
      <RideList {...query} />
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

const ColumnView = styled.div`
  display: flex;
  flex-direction: row;
  gap: 2rem;
`;
