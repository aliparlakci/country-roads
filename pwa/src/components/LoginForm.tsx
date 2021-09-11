import React, { useState } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import CONSTANTS from '../constants'
import { useHistory } from 'react-router-dom'

export default function LoginForm() {
  const history = useHistory()
  const [loading, setLoading] = useState(false)
  const [email, setEmail] = useState("")

  const handleSubmit: React.FormEventHandler<HTMLFormElement> = async (
    event,
  ) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)

    let response
    try {
      response = await fetch(CONSTANTS.API.AUTH.LOGIN, {
        method: 'POST',
        body: formData,
      })
    } catch (e) {
      console.error(e)
      return
    }

    if (response.status === 200) {
      history.push(`${CONSTANTS.ROUTES.OTP}?email=${email}&redirect=${CONSTANTS.ROUTES.RIDES.MAIN}`)
      return
    }

    if (response.status === 201) {
      history.push(`${CONSTANTS.ROUTES.OTP}?email=${email}&redirect=${CONSTANTS.ROUTES.SETTINGS}`)
      return
    }

    setLoading(false)
  }

  return (
    <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label
              htmlFor="email"
              className="block text-base sm:text-sm font-medium text-gray-700"
            >
              Email address
            </label>
            <div className="mt-1">
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                pattern="^.+@sabanciuniv\.edu$"
                required
                autoFocus
                onChange={(event) => setEmail(event.target.value)}
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <span className="text-center text-xs text-gray-400">
            Use your @sabanciuniv.edu email address
          </span>

          {!loading && (
            <div>
              <button
                type="submit"
                className="w-full flex items-center gap-2 justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Continue
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
        </form>
      </div>
    </div>
  )
}
