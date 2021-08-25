import React from "react";
import styled from "styled-components";
import IRide from "../models/ride";
import capitalize from "../utils/capitalize";

export interface IRideItemProps {
  ride: IRide;
}

export default function RideItem({ ride }: IRideItemProps) {
  return (
    <RideItemContainer>
      <b>{capitalize(ride.type)}</b>
      <span>
        From{" "}
        {ride.direction === "from_campus" ? "Campus" : ride.destination.display}
      </span>
      <span>
        To{" "}
        {ride.direction === "to_campus" ? "Campus" : ride.destination.display}
      </span>
      <span>On {new Date(ride.date * 1000).toDateString()}</span>
    </RideItemContainer>
  );
}

const RideItemContainer = styled.div`
  display: flex;
  flex-direction: row;
  gap: 0.5rem;
`;
