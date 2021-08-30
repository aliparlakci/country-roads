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
      {locationResponse && locationResponse.results && (
        <>
          <option value="none">None</option>
          {locationResponse?.results?.map((location) => (
            <option key={location.key} value={location.key}>
              {location.display}
            </option>
          ))}
        </>
      )}
      {!locationResponse && <option disabled={true}>Loading...</option>}
      {locationResponse && locationResponse.error && (
        <option disabled={true}>Error!</option>
      )}
    </select>
  );
}
