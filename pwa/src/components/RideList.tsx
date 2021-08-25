import React from "react";
import useSWR from "swr";

import RideItem from "./RideItem";
import CONSTANTS from "../constants";
import IRide from "../models/ride";

export interface IRideListProps {
  refresh: number
}

export default function RideList(props: IRideListProps) {
  const URI = process.env.BACKEND_URI || "http://localhost:8080";

  const { data, error } = useSWR<IRide[]>(`${URI}/${CONSTANTS.API.RIDES}`);

  if (!data) return <div>Loading...</div>;
  if (error) return <div>Error</div>;

  return (
    <ol>
      {data.map((ride, i) => (
        <li key={i}>
          <RideItem ride={ride} />
        </li>
      ))}
    </ol>
  );
}
