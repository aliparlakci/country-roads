import React from "react";
import useSWR from "swr";

import ILocation from "../models/location";
import CONSTANTS from "../constants";

export default function useLocations() {
  const { data, error } = useSWR<ILocation[]>(CONSTANTS.API().RIDES);
  if (error) console.error(error);
  return data;
}
