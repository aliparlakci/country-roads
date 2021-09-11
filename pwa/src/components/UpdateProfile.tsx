import React, { useCallback, useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { mutate } from 'swr'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

import CONSTANTS from '../constants'
import IContactInfo from '../types/contact'
import useAuth from '../hooks/useAuth'
import { IUser } from '../types/user'

export default function UpdateProfile() {
  const [loading, setLoading] = React.useState(false)
  const [previousContact, setPreviousContact] = useState<IContactInfo>({})

  const { user } = useAuth()
  const history = useHistory()

  const fetchDefault = async (user: IUser) => {
    try {
      const reponse = await fetch(CONSTANTS.API.CONTACT(user.id))
      setPreviousContact(await reponse.json())
    } catch (e) {
      console.log(e)
    }
  }

  useEffect(() => {
    user && fetchDefault(user)
  }, [user])

  const handleSubmit = async (event: any) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)

    let result
    try {
      result = await fetch(CONSTANTS.API.CONTACT(''), {
        method: 'PUT',
        body: formData,
      })
    } catch (e) {
      console.error(e)
    }

    if (!result || result.status !== 200) {
      console.log(result?.status)
      setLoading(false)
      return
    }

    await mutate(CONSTANTS.API.USERS.MAIN)

    event.target.reset()
    history.push(CONSTANTS.ROUTES.RIDES.MAIN)
  }

  return (
    <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div className="bg-white py-4 sm:py-8 px-4 shadow rounded-lg sm:px-10">
        <div className="text-2xl font-semibold mb-4">Update contact info</div>
        <form className="space-y-6" onSubmit={handleSubmit}>
          <div>
            <label
              htmlFor="name"
              className="block text-base sm:text-sm font-medium text-gray-700"
            >
              Name
            </label>
            <div className="mt-1">
              <input
                id="name"
                name="displayName"
                type="name"
                autoComplete="name"
                defaultValue={previousContact.name}
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label
              htmlFor="email"
              className="block text-base sm:text-sm font-medium text-gray-700"
            >
              Phone
            </label>
            <div className="mt-1">
              <input
                id="phone"
                name="phone"
                type="tel"
                autoComplete="tel"
                pattern="^(\+|)([0-9]{1,3})([0-9]{10})$"
                defaultValue={previousContact.phone}
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label
              htmlFor="email"
              className="block text-base sm:text-sm font-medium text-gray-700"
            >
              Whatsapp
            </label>
            <div className="mt-1">
              <input
                id="whatsapp"
                name="whatsapp"
                type="tel"
                autoComplete="tel"
                pattern="^(\+|)([0-9]{1,3})([0-9]{10})$"
                defaultValue={previousContact.whatsapp}
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <p className="text-center text-sm sm:text-xs text-gray-400">
            Leave empty the fields which you do not want to specify. Other users
            can still get in touch with you through email.
          </p>

          <div className="flex flex-row gap-2">
            {!loading && (
              <button
                type="submit"
                className="w-full flex items-center gap-2 justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Update
              </button>
            )}
            {loading && (
              <div className="flex justify-center w-full">
                <FontAwesomeIcon
                  className="animate-spin text-lg text-indigo-600"
                  icon={faCircleNotch}
                />
              </div>
            )}
            <button
              type="reset"
              onClick={() => user && fetchDefault(user)}
              className="w-full flex items-center gap-2 justify-center py-2 px-4 border-2 rounded-md shadow-sm text-base sm:text-sm font-medium text-indigo-600 border-indigo-600 hover:border-indigo-800 hover:text-indigo-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Revert
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
