import React from 'react'
import styled from 'styled-components'
import CONSTANTS from '../constants'
import useAuth from '../hooks/useAuth'
import IRide from '../types/ride'
import capitalize from '../utils/capitalize'
import mutateWithQueries from '../utils/mutateWithQueries'

export interface IRideItemProps {
  ride: IRide
}

export default function RideItem({ ride }: IRideItemProps) {
  const { user } = useAuth()

  const doDelete = async () => {
    try {
      await fetch(CONSTANTS.API.RIDE(ride.id), { method: 'delete' })
      mutateWithQueries(CONSTANTS.API.RIDES)
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <RideItemContainer>
      <TitleArea>
        <b>{capitalize(ride.type)}</b>
        {ride.owner.id === user?.id && (
          <CloseButton onClick={doDelete}>
            <b>X</b>
          </CloseButton>
        )}
      </TitleArea>
      <span>
        From{' '}
        {ride.direction === 'from_campus' ? 'Campus' : ride.destination.display}
      </span>
      <span>
        To{' '}
        {ride.direction === 'to_campus' ? 'Campus' : ride.destination.display}
      </span>
      <i>on {new Date(ride.date * 1000).toDateString()}</i>
      <span>{ride.owner.displayName}</span>
    </RideItemContainer>
  )
}

const RideItemContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.5rem;
  border-radius: 0.5rem;
  border: solid 1px lightgray;

  transition: 0.1s;

  & > :not(:first-child) {
    text-align: right;
  }
`

const TitleArea = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
`

const CloseButton = styled.button`
  border: none;
  outline: none;
  background: none;
  color: red;

  &:hover {
    transform: translateX(-0.1rem) translateY(-0.1rem);
    transition: 100ms;
  }
`
