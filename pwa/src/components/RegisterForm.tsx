import React, { useState } from "react";
import styled from "styled-components";
import { mutate } from "swr";

import CONSTANTS from "../constants";

export interface IRegisterFormProps {}

export default function RegisterForm(props: IRegisterFormProps) {
  const [disabled, setDisabled] = useState(false);

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    setDisabled(true);

    const formData = new FormData(event.currentTarget);

    try {
      await fetch(CONSTANTS.API().USERS, {
        method: "POST",
        body: formData,
      });
    } catch (e) {
      console.error(e);
    }

    mutate(CONSTANTS.API().USERS);

    event.target.reset();
    setDisabled(false);
  };

  return (
    <StyledForm onSubmit={handleSubmit}>
      <FormItem>
        <span>Name:</span>
        <input
          type="text"
          name="displayName"
          placeholder="Enter your name..."
          required
        />
      </FormItem>

      <FormItem>
        <span>Email (@sabanciuniv email):</span>
        <input
          type="email"
          name="email"
          placeholder="@sabanciuniv.edu"
          pattern="^.+@sabanciuniv\.edu$"
          required
        />
      </FormItem>

      <FormItem>
        <span>Phone:</span>
        <div>
          <input
            type="tel"
            name="phone"
            pattern="^(\+|)([0-9]{1,3})([0-9]{10})$"
            required
          />
        </div>
      </FormItem>

      <input type="submit" value="Register" disabled={disabled} />
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

const FormItem = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: start;
  gap: 0.25rem;
`;
