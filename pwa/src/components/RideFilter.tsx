import React, { useEffect, useReducer } from "react";
import styled from "styled-components";
import LocationsDropdown from "./LocationsDropdown";
import { useHistory, useLocation } from "react-router-dom";

type RideFilterReducer = (
  state: IRideFilterState,
  action: IRideFilterReducerAction
) => IRideFilterState;

interface IRideFilterReducerAction {
  type: string;
  field: string;
  payload: string | boolean;
}

interface IRideFilterState {
  [key: string]: { value: string; disabled: boolean };

  type: { value: string; disabled: boolean };
  direction: { value: string; disabled: boolean };
  destination: { value: string; disabled: boolean };
  startDate: { value: string; disabled: boolean };
  endDate: { value: string; disabled: boolean };
}

const initialState = {
  type: { value: "", disabled: true },
  direction: { value: "", disabled: true },
  destination: { value: "", disabled: true },
  startDate: { value: "", disabled: true },
  endDate: { value: "", disabled: true },
};

const reducer: RideFilterReducer = (state, action) => {
  switch (action.type) {
    case "set":
      return {
        ...state,
        [action.field]: { ...state[action.field], value: action.payload },
      } as IRideFilterState;
    case "disable":
      return {
        ...state,
        [action.field]: { ...state[action.field], disabled: action.payload },
      } as IRideFilterState;
    default:
      return state;
  }
};

export default function RideFilter() {
  const [state, dispatch] = useReducer(reducer, initialState);
  const history = useHistory();
  const { pathname } = useLocation();

  useEffect(() => {
    const params = new URLSearchParams();
    if (!state.type.disabled && state.type.value)
      params.set("type", state.type.value);
    if (!state.destination.disabled && state.destination.value)
      params.set("destination", state.destination.value);
    if (!state.direction.disabled && state.direction.value)
      params.set("direction", state.direction.value);
    if (!state.startDate.disabled && state.startDate.value)
      params.set("start_date", state.startDate.value);
    if (!state.startDate.disabled && state.endDate.value)
      params.set("end_date", state.endDate.value);

    history.push(`${pathname}?${params.toString()}`);
  }, [state, history, pathname]);

  return (
    <Container>
      <ParameterContainer>
        <LabelContainer>
          <input
            type="checkbox"
            id="disableType"
            onChange={(event) =>
              dispatch({
                type: "disable",
                field: "type",
                payload: !event.target.checked,
              })
            }
          />
          <label htmlFor="disableType">Type</label>
        </LabelContainer>
        <select
          onChange={(event) =>
            dispatch({
              field: "type",
              type: "set",
              payload: event.target.value,
            })
          }
          disabled={state.type.disabled}
        >
          <option value=""></option>
          <option value="request">Request</option>
          <option value="offer">Offer</option>
          <option value="taxi">Taxi</option>
        </select>
      </ParameterContainer>
      <ParameterContainer>
        <LabelContainer>
          <input
            type="checkbox"
            id="disableDirection"
            onChange={(event) =>
              dispatch({
                type: "disable",
                field: "direction",
                payload: !event.target.checked,
              })
            }
          />
          <label htmlFor="disableDirection">Direction</label>
        </LabelContainer>
        <select
          onChange={(event) =>
            dispatch({
              field: "direction",
              type: "set",
              payload: event.target.value,
            })
          }
          disabled={state.direction.disabled}
        >
          <option value=""></option>
          <option value="to_campus">To Campus</option>
          <option value="from_campus">From Campus</option>
        </select>
      </ParameterContainer>
      <ParameterContainer>
        <LabelContainer>
          <input
            type="checkbox"
            id="disableDestination"
            onChange={(event) =>
              dispatch({
                type: "disable",
                field: "destination",
                payload: !event.target.checked,
              })
            }
          />
          <label htmlFor="disableDestination">Destination</label>
        </LabelContainer>
        <LocationsDropdown
          onChange={(event) =>
            dispatch({
              field: "destination",
              type: "set",
              payload: event.target.value,
            })
          }
          disabled={state.destination.disabled}
          onData={() => null}
        />
      </ParameterContainer>
      <ParameterContainer>
        <input
          type="checkbox"
          id="disableStartDate"
          onChange={(event) =>
            dispatch({
              type: "disable",
              field: "startDate",
              payload: !event.target.checked,
            })
          }
        />
        <input
          type="date"
          disabled={state.startDate.disabled}
          onChange={(event) =>
            dispatch({
              type: "set",
              field: "startDate",
              payload: (
                new Date(event.target.value).getTime() / 1000
              ).toString(),
            })
          }
        />
        <input
          type="date"
          disabled={state.endDate.disabled}
          onChange={(event) =>
            dispatch({
              type: "set",
              field: "endDate",
              payload: (
                new Date(event.target.value).getTime() / 1000
              ).toString(),
            })
          }
        />
        <input
          type="checkbox"
          id="disableEndDate"
          onChange={(event) =>
            dispatch({
              type: "disable",
              field: "endDate",
              payload: !event.target.checked,
            })
          }
        />
      </ParameterContainer>
    </Container>
  );
}

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;

  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid lightgray;
`;

const ParameterContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 0.5rem;
`;

const LabelContainer = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
`;
