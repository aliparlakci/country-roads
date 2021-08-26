import React, { useEffect, useState } from "react";
import styled from "styled-components";

import CONSTANTS from "../constants";
import useLocations from "../hooks/useLocations";
import mutateWithQueries from "../utils/mutateWithQueries";

export interface INewRideFormProps {}

export default function NewRideForm(props: INewRideFormProps) {
  const { data: locations } = useLocations();
  const [disabled, setDisabled] = useState(false);

  useEffect(() => setDisabled(!locations), [locations]);

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    setDisabled(true);

    const formData = new FormData(event.currentTarget);
    const date = formData.get("date")?.toString();
    if (date)
      formData.set("date", (new Date(date).getTime() / 1000).toString());

    try {
      await fetch(CONSTANTS.API().RIDES, {
        method: "POST",
        body: formData,
      });
    } catch (e) {
      console.error(e);
    }

    mutateWithQueries(CONSTANTS.API().RIDES)

    event.target.reset();
    setDisabled(false);
  };

  return (
    <StyledForm onSubmit={handleSubmit}>
      <FormItem>
        <span>Ride type:</span>
        <StyledUl>
          <li>
            <input
              type="radio"
              name="type"
              id="type_request"
              value="request"
              required
              disabled={disabled}
            />
            <label htmlFor="type_request">Request a ride</label>
          </li>
          <li>
            <input
              type="radio"
              name="type"
              id="type_offer"
              value="offer"
              required
              disabled={disabled}
            />
            <label htmlFor="type_offer">Offer a ride</label>
          </li>
          <li>
            <input
              type="radio"
              name="type"
              id="type_taxi"
              value="taxi"
              required
              disabled={disabled}
            />
            <label htmlFor="type_taxi">Share a taxi</label>
          </li>
        </StyledUl>
      </FormItem>

      <FormItem>
        <div>
          <input
            type="radio"
            name="direction"
            id="direction_to"
            value="to_campus"
            required
            disabled={disabled}
          />
          <label htmlFor="direction_to">To campus</label>
        </div>
        <div>
          <input
            type="radio"
            name="direction"
            id="direction_from"
            value="from_campus"
            required
            disabled={disabled}
          />
          <label htmlFor="direction_from">From campus</label>
        </div>
      </FormItem>

      <FormItem>
        <label htmlFor="destination">Destination:</label>
        <select
          id="destination"
          name="destination"
          required
          disabled={disabled}
        >
          {locations &&
            locations.map((location) => (
              <option key={location.id} value={location.id}>
                {location.display}
              </option>
            ))}
          {!locations && <option>Loading...</option>}
        </select>
      </FormItem>

      <FormItem>
        <label htmlFor="name">When:</label>
        <input
          type="date"
          id="name"
          name="date"
          min={new Date().toISOString().substring(0, 10)} // Today's date string
          required
          disabled={disabled}
        />
      </FormItem>

      <input type="submit" value="Post" disabled={disabled} />
    </StyledForm>
  );
}

const StyledForm = styled.form`
  width: 16rem;

  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 1rem;

  border-radius: 1rem;
  border: 1px solid lightgrey;
`;

const StyledUl = styled.ul`
  padding: 0;
  margin: 0;
  list-style: none;
`;

const FormItem = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`;
