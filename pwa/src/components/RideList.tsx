import React from 'react'

import RideItem from './RideItem'
import useRides, { IRideQuery } from '../hooks/useRides'

export interface IRideListProps extends IRideQuery {}

export default function RideList(props: IRideListProps) {
  const { data, error } = useRides(props)

  if (!data) return <i>Loading...</i>
  if (data && !data.results) return <i>Nothing to see here...</i>
  if (error) return <div>Error</div>
  if (data && data.error) return <div>Error</div>

  return (
    <div className="grid gap-2 grid-cols-1 sm:grid-cols-3">
      {data &&
        data.results &&
        data.results.map((ride, i) => <RideItem key={i} ride={ride} />)}
    </div>
  )
}
