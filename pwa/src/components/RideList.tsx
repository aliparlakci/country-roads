import React from "react";

import RideItem from "./RideItem";
import styled from "styled-components";
import useRides, { IRideQuery } from "../hooks/useRides";

export interface IRideListProps extends IRideQuery {}

export default function RideList(props: IRideListProps) {
  const { data, error } = useRides(props);

  if (!data) return <i>Nothing to see here...</i>;
  if (error) return <div>Error</div>;
  if (data && data.error) return <div>Error</div>;

  return (
    <RideListContainer>
      {data &&
        data.results &&
        data.results.map((ride, i) => <RideItem key={i} ride={ride} />)}
    </RideListContainer>
  );
}

const RideListContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  grid-template-rows: auto;
  flex-direction: column;
  gap: 0.5rem;
  transition: 0.1s;
`;
