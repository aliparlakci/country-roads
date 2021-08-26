import useSWR from "swr";

import ILocation from "../types/location";
import CONSTANTS from "../constants";

export default function useLocations() {
  const { data, error } = useSWR<ILocation[]>(CONSTANTS.API().LOCATIONS);
  if (error) console.error(error);
  return { data, error };
}
