import useSWR from "swr";

import IRide from "../types/ride";
import CONSTANTS from "../constants";

export default function useRides() {
  const { data, error } = useSWR<IRide[]>(CONSTANTS.API().RIDES);
  if (error) console.error(error);
  return { data, error };
}
