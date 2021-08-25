import React from "react";
import useSWR from "swr";
import CONSTANTS from "../constants";

export default function RideList() {
  const URI = process.env.BACKEND_URI || "http://localhost:8080";

  const { data, error } = useSWR(`${URI}/${CONSTANTS.API.RIDES}`);

  if (!data)
    return (
      <div>
        <h2>Ride List</h2>
      </div>
    );

  if (error) return <div>Error</div>;

  return <div>{JSON.stringify(data)}</div>;
}
