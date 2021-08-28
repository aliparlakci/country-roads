import useSWR from "swr";

import IRide from "../types/ride";
import CONSTANTS from "../constants";

export interface IRideResponse {
  results?: IRide[];
  error?: string
}

export default function useRides() {
  const { data, error } = useSWR<IRideResponse>(CONSTANTS.API().RIDES);
  if (error) console.error(error);
  return { data, error };
}
