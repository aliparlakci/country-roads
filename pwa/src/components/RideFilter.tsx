import React, { useEffect, useState } from 'react'
import { useHistory, useLocation } from 'react-router'
import LocationsDropdown from './LocationsDropdown'

export default function RideFilter() {
  const [type, setType] = useState('')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')

  const history = useHistory()
  const { pathname } = useLocation()

  useEffect(() => {
    const params = new URLSearchParams()
    if (type) params.set('type', type)
    if (from) params.set('from', from)
    if (to) params.set('to', to)
    if (startDate) params.set('start_date', startDate)
    if (endDate) params.set('end_date', endDate)

    history.push(`${pathname}?${params.toString()}`)
  }, [type, from, to, startDate, endDate, history, pathname])
  return (
    <div>
      <div
        className="grid grid-cols-2 sm:flex flex-row justify-center items-center gap-2 w-full px-3 py-3 rounded-t-lg bg-white border-b border-gray-200 shadow"
        style={{
          gridTemplateColumns: 'auto 1fr',
        }}
      >
        <label
          htmlFor="type_filter"
          className="text-base sm:text-sm text-right"
        >
          Type
        </label>
        <select
          id="type_filter"
          className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
          onChange={(event) => setType(event.target.value)}
        >
          <option>All</option>
          <option value="request">Ride request</option>
          <option value="offer">Ride offer</option>
          <option value="taxi">Share a taxi</option>
        </select>
        <label
          htmlFor="from_filter"
          className="text-base sm:text-sm text-right"
        >
          From
        </label>
        <LocationsDropdown
          id="from_filter"
          className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
          onChange={(event) => setFrom(event.target.value)}
        >
          <option value="">All</option>
        </LocationsDropdown>
        <label htmlFor="to_filter" className="text-base sm:text-sm text-right">
          To
        </label>
        <LocationsDropdown
          id="to_filter"
          className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
          onChange={(event) => setTo(event.target.value)}
        >
          <option value="">All</option>
        </LocationsDropdown>
      </div>
      <div
        className="grid grid-cols-2 sm:flex flex-row justify-center items-center gap-2 w-full px-3 py-3 rounded-b-lg bg-white shadow"
        style={{
          gridTemplateColumns: 'auto auto',
        }}
      >
        <label
          htmlFor="start_date_filter"
          className="text-base sm:text-sm text-right"
        >
          Start date
        </label>
        <input
          type="date"
          id="start_date_filter"
          className="block text-base sm:text-sm border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 rounded-md"
          onChange={(event) =>
            setStartDate(
              (new Date(event.target.value).getTime() / 1000).toString(),
            )
          }
        />
        <label
          htmlFor="end_date_filter"
          className="text-base sm:text-sm text-right"
        >
          End date
        </label>
        <input
          type="date"
          id="end_date_filter"
          className="block text-base sm:text-sm border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 rounded-md"
          onChange={(event) =>
            setEndDate(
              (new Date(event.target.value).getTime() / 1000).toString(),
            )
          }
        />
      </div>
    </div>
  )
}
