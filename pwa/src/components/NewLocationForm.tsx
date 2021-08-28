import React, { useState } from "react";
import styled from "styled-components";
import { mutate } from "swr";

import CONSTANTS from "../constants";
import LocationsDropdown from "./LocationsDropdown";

export interface INewLocationFormProps {}

export default function NewLocationForm(props: INewLocationFormProps) {
  const [disabled, setDisabled] = useState(true);

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    setDisabled(true);

    const formData = new FormData(event.currentTarget);
    if (formData.get("parentId")?.valueOf() === "-1")
      formData.delete("parentId");

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
        <LocationsDropdown
          id="parentId"
          name="parentId"
          required
          disabled={disabled}
          onData={() => setDisabled(false)}
        />
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
