import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { useState } from 'react'
import { useHistory } from 'react-router-dom'
import { mutate } from 'swr'
import CONSTANTS from '../constants'

export default function OTPForm({ email }: { email: string }) {
  const [loading, setLoading] = useState(false)
  const history = useHistory()

  const handleLogin = async (event: any) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)
    formData.set('email', email)
    let result
    try {
      result = await fetch(CONSTANTS.API.AUTH.VERIFY, {
        method: 'POST',
        body: formData,
      })
    } catch (e) {
      console.error(e)
    }

    if (!result || result.status !== 200) {
      setLoading(false)
      console.log({ error: await result?.body?.getReader().read() })
      return
    }

    mutate(CONSTANTS.API.AUTH.USER)
    history.push(CONSTANTS.ROUTES.RIDES.MAIN)
  }

  return (
    <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form className="space-y-6" onSubmit={handleLogin}>
          <div>
            <label
              htmlFor="otp"
              className="block text-base sm:text-sm font-medium text-gray-700"
            >
              Your one time password
            </label>
            <div className="mt-1">
              <input
                id="otp"
                name="otp"
                type="numeric"
                autoComplete="one-time-password"
                required
                autoFocus
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <p className="text-sm sm:text-xs text-gray-400 text-center">
            We have sent your email a one time password. Enter it to login to
            SUPool
          </p>
          <div>
            {!loading && (
              <button
                type="submit"
                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Sign in
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
          </div>
        </form>
      </div>
    </div>
  )
}
