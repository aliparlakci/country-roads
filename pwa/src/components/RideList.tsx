import React from "react";
import useSWR from "swr";

import RideItem from "./RideItem";
import CONSTANTS from "../constants";
import IRide from "../models/ride";
import styled from "styled-components";

export interface IRideListProps {
  refresh: number;
}

export default function RideList(props: IRideListProps) {
  const { data, error } = useSWR<IRide[]>(CONSTANTS.API().RIDES);

  if (!data) return <i>Nothing to see here...</i>;
  if (error) return <div>Error</div>;

  return (
    <RideListContainer>
      {data.map((ride, i) => (
        <RideItem key={i} ride={ride} />
      ))}
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
