import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React from 'react'
import { useHistory } from 'react-router-dom'
import { mutate } from 'swr'
import CONSTANTS from '../constants'

export default function RegisterForm({ email }: { email?: string }) {
  const [loading, setLoading] = React.useState(false)
  const history = useHistory()

  const handleSubmit = async (event: any) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)

    let result
    try {
      result = await fetch(CONSTANTS.API.USERS, {
        method: 'POST',
        body: formData,
      })
    } catch (e) {
      console.error(e)
    }

    if (!result || result.status !== 201) {
      console.log(result?.status)
      setLoading(false)
      return
    }

    await mutate(CONSTANTS.API.USERS)
    
    event.target.reset()
    history.push(`${CONSTANTS.ROUTES.LOGIN}?email=${formData.get('email')}`)
  }

  return (
    <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <form className="space-y-6" onSubmit={handleSubmit}>
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
                defaultValue={email}
                required
                autoFocus
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>
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
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
                required
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
                required
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          {!loading && (
            <div>
              <button
                type="submit"
                className="w-full flex items-center gap-2 justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Register
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
