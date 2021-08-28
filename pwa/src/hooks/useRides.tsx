import useSWR from "swr";

import IRide from "../types/ride";
import CONSTANTS from "../constants";
import filterBuilder from "../utils/filterBuilder";

export interface IRideResponse {
  results?: IRide[];
  error?: string;
}

export interface IRideQuery {
  type: string | null;
  direction: string | null;
  destination: string | null;
  startDate: string | null;
  endDate: string | null;
}

export default function useRides(query: IRideQuery) {
  const endpoint = `${CONSTANTS.API().RIDES}?${filterBuilder(query)}`;
  console.log(endpoint);
  const { data, error } = useSWR<IRideResponse>(endpoint);
  if (error) console.error(error);
  return { data, error };
}
