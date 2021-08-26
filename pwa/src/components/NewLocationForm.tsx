import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { mutate } from "swr";

import CONSTANTS from "../constants";
import useLocations from "../hooks/useLocations";

export interface INewLocationFormProps {}

export default function NewLocationForm(props: INewLocationFormProps) {
  const { data: locations } = useLocations();
  const [disabled, setDisabled] = useState(false);

  useEffect(() => setDisabled(!locations), [locations]);

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    setDisabled(true);

    const formData = new FormData(event.currentTarget);
    if (formData.get("parentId")?.valueOf() === -1) formData.delete("parentId");

    try {
      await fetch(CONSTANTS.API().LOCATIONS, {
        method: "POST",
        body: formData,
      });
    } catch (e) {
      console.error(e);
    }

    mutate(CONSTANTS.API().LOCATIONS);

    event.target.reset();
    setDisabled(false);
  };

  return (
    <StyledForm onSubmit={handleSubmit}>
      <FormItem>
        <label htmlFor="displayName">Display name:</label>
        <input
          type="text"
          placeholder="Type display name..."
          id="displayName"
          name="display"
          disabled={disabled}
        />
      </FormItem>

      <FormItem>
        <label htmlFor="destination">Parent:</label>
        <select id="parentId" name="parentId" required disabled={disabled}>
          <option value={-1}>None</option>
          {locations &&
            locations.map((location) => (
              <option key={location.id} value={location.id}>
                {location.display}
              </option>
            ))}
          {!locations && <option>Loading...</option>}
        </select>
      </FormItem>

      <input type="submit" value="Post" disabled={disabled} />
    </StyledForm>
  );
}

const StyledForm = styled.form`
  width: fit-content;

  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 1rem;

  border-radius: 1rem;
  border: 1px solid lightgrey;
`;

const FormItem = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
`;
