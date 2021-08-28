import React, { useEffect } from "react";
import useLocations from "../hooks/useLocations";

export interface ILocationsDropdownProps
  extends React.DetailedHTMLProps<
    React.SelectHTMLAttributes<HTMLSelectElement>,
    HTMLSelectElement
  > {
  onData: CallableFunction;
}

export default function LocationsDropdown(props: ILocationsDropdownProps) {
  const { onData } = props;
  const { data: locationResponse } = useLocations();

  useEffect(() => locationResponse && onData(), [onData, locationResponse]);

  return (
    <select {...props}>
      <option value={-1}>None</option>
      {locationResponse && locationResponse.results && (
        <>
          {locationResponse?.results?.map((location) => (
            <option key={location.key} value={location.key}>
              {location.display}
            </option>
          ))}
        </>
      )}
      {!locationResponse && <option>Loading...</option>}
      {locationResponse && locationResponse.error && <option>Error!</option>}
    </select>
  );
}
