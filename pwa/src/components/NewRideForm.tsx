import React, { useState } from 'react'
import cn from 'classnames'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

import CONSTANTS from '../constants'
import mutateWithQueries from '../utils/mutateWithQueries'
import { useHistory } from 'react-router-dom'
import { SportsCar, Student, Taxi } from './icons'
import LocationsDropdown from './LocationsDropdown'

export interface INewRideFormProps {}

export default function NewRideForm(props: INewRideFormProps) {
  const [loading, setLoading] = useState(false)
  const [selectedType, setSelectedType] = useState('taxi')
  const history = useHistory()

  const handleSubmit = async (event: any) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)
    const date = formData.get('date')?.toString()
    if (date) formData.set('date', (new Date(date).getTime() / 1000).toString())

    try {
      await fetch(CONSTANTS.API.RIDES, {
        method: 'POST',
        body: formData,
      })
    } catch (e) {
      console.error(e)
    }

    mutateWithQueries(CONSTANTS.API.RIDES)

    event.target.reset()
    setLoading(false)

    history.push(CONSTANTS.ROUTES.RIDES.MAIN)
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="mt-8 sm:mx-auto sm:w-full sm:max-w-md"
    >
      <div className="flex flex-row rounded-t-lg shadow bg-white border-b border-gray-100">
        <button
          onClick={(e) => {
            e.preventDefault()
            setSelectedType('offer')
          }}
          className={cn(
            'border-b-2 hover:border-indigo-600 border-transparent rounded-tl-lg flex flex-col justify-center items-center px-2 py-2 w-full',
            { 'border-indigo-600': selectedType === 'offer' },
          )}
        >
          <SportsCar className="text-5xl" />
          <span className="text-base sm:text-sm text-center">Offer a lift</span>
        </button>
        <button
          onClick={(e) => {
            e.preventDefault()
            setSelectedType('taxi')
          }}
          className={cn(
            'border-b-2 hover:border-indigo-600 border-transparent flex flex-col justify-center items-center px-2 py-2 w-full',
            { 'border-indigo-600': selectedType === 'taxi' },
          )}
        >
          <Taxi className="text-5xl" />
          <span className="text-base sm:text-sm text-center">Share a taxi</span>
        </button>
        <button
          onClick={(e) => {
            e.preventDefault()
            setSelectedType('request')
          }}
          className={cn(
            'border-b-2 hover:border-indigo-600 border-transparent rounded-tr-lg flex flex-col justify-center items-center px-2 py-2 w-full',
            { 'border-indigo-600': selectedType === 'request' },
          )}
        >
          <Student className="text-5xl" />
          <span className="text-base sm:text-sm text-center">
            Request a lift
          </span>
        </button>
        <input type="hidden" name="type" value={selectedType} />
      </div>
      <div className="flex flex-col gap-3 bg-white py-4 px-4 shadow rounded-b-lg sm:px-10">
        <div>
          <label
            htmlFor="from"
            className="block text-base sm:text-sm font-medium text-gray-700"
          >
            From
          </label>
          <div className="mt-1">
            <LocationsDropdown
              id="from"
              name="from"
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            />
          </div>
        </div>
        <div>
          <label
            htmlFor="to"
            className="block text-base sm:text-sm font-medium text-gray-700"
          >
            To
          </label>
          <div className="mt-1">
            <LocationsDropdown
              id="to"
              name="to"
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="date"
            className="block text-base sm:text-sm font-medium text-gray-700"
          >
            When
          </label>
          <div className="mt-1">
            <input
              type="date"
              id="date"
              name="date"
              required
              className="mt-1 block w-full pl-3 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            />
          </div>
        </div>

        {!loading && (
          <div>
            <button
              type="submit"
              className="w-full flex items-center gap-2 justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Create
            </button>
          </div>
        )}
        {loading && (
          <div className="flex justify-center w-full">
            <FontAwesomeIcon
              className="animate-spin text-lg text-indigo-600"
              icon={faCircleNotch}
            />
          </div>
        )}
      </div>
    </form>
  )
}
