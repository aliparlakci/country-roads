import React, {useReducer} from 'react'
import styled from "styled-components";
import LocationsDropdown from "./LocationsDropdown";

type RideFilterReducer = (state: IRideFilterState, action: IRideFilterReducerAction) => IRideFilterState

interface IRideFilterReducerAction {
  type: string,
  field: string,
  payload: string | boolean,
}

interface IRideFilterState {
  [key: string]: { value: string, disabled: boolean }

  type: { value: string, disabled: boolean },
  direction: { value: string, disabled: boolean },
  destination: { value: string, disabled: boolean },
  startDate: { value: string, disabled: boolean },
  endDate: { value: string, disabled: boolean }
}

const initialState = {
  type: {value: "", disabled: true},
  direction: {value: "", disabled: true},
  destination: {value: "", disabled: true},
  startDate: {value: "", disabled: true},
  endDate: {value: "", disabled: true},
}

const reducer: RideFilterReducer = (state, action) => {
  switch (action.type) {
    case 'value':
      return {...state, [action.field]: {...state[action.field], value: action.payload}} as IRideFilterState
    case 'disable':
      return {...state, [action.field]: {...state[action.field], disabled: action.payload}} as IRideFilterState
    default:
      return state
  }
}

export default function RideFilter() {
  const [state, dispatch] = useReducer(reducer, initialState)

  return <Container>
    <ParameterContainer>
      <Label>
        <input type="checkbox"
               onChange={(event) => dispatch({type: "disable", field: "type", payload: !event.target.checked})}/>
        <span>Type</span>
      </Label>
      <select disabled={state.type.disabled}>
        <option></option>
        <option>Request</option>
        <option>Offer</option>
        <option>Taxi</option>
      </select>
    </ParameterContainer>
    <ParameterContainer>
      <Label>
        <input type="checkbox"/>
        <span>Direction</span>
      </Label>
      <select>
        <option></option>
        <option>To Campus</option>
        <option>From Campus</option>
      </select>
    </ParameterContainer>
    <ParameterContainer>
      <Label>
        <input type="checkbox"/>
        <span>Direction</span>
      </Label>
      <LocationsDropdown onData={() => null}/>
    </ParameterContainer>
    <ParameterContainer>
      <input type="date" value={(new Date()).getTime()}/><input type="date"/>
    </ParameterContainer>
  </Container>
}

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;

  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid lightgray;
`

const ParameterContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 0.5rem;
`

const Label = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
`