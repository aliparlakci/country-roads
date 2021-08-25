import React from "react";
import useSWR from "swr";

import ILocation from "../models/location";
import CONSTANTS from "../constants";

export default function useLocations() {
  const URI = process.env.BACKEND_URI || "http://localhost:8080";

  const { data, error } = useSWR<ILocation[]>(
    `${URI}/${CONSTANTS.API.LOCATIONS}`
  );

  if (error) console.error(error);

  return data;
}
