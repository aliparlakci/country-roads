import useSWR from "swr";

import IRide from "../types/ride";
import CONSTANTS from "../constants";
import RideType from "../types/rideType";
import RideDirection from "../types/rideDirection";

export interface IRideResponse {
  results?: IRide[];
  error?: string;
}

export interface IRideQuery {
  type?: RideType;
  direction?: RideDirection;
  startDate?: Date;
  endDate?: Date;
}

export default function useRides({type, direction, startDate, endDate}: IRideQuery) {
  const endpoint = `${CONSTANTS.API().RIDES}?type=${type}&direction=${direction}&start_date=${startDate && (startDate.getTime() / 1000)}&end_date=${endDate && (endDate?.getTime()) / 1000}`

  const { data, error } = useSWR<IRideResponse>(endpoint);
  if (error) console.error(error);
  return { data, error };
}
