import React from 'react'
import { useHistory } from 'react-router-dom'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import useAuth from '../hooks/useAuth'
import CONSTANTS from '../constants'
import { mutate } from 'swr'
import useModal from '../hooks/useModal'

export default function UpdateProfile({ email }: { email?: string }) {
  const [loading, setLoading] = React.useState(false)
  const { user } = useAuth()
  const { alert, error } = useModal()

  const handleSubmit = async (event: any) => {
    event.preventDefault()
    setLoading(true)

    const formData = new FormData(event.currentTarget)

    let result
    try {
      result = await fetch(`${CONSTANTS.API.USERS.MAIN}/${user?.id}`, {
        method: 'PUT',
        body: formData,
      })
    } catch (e) {
      console.error(e)
    }

    if (!result || result.status !== 200) {
      console.log(result?.status)
      setLoading(false)
      error({ header: 'Error', body: 'Profile cannot get updated' })
      return
    }

    await mutate(CONSTANTS.API.USERS.MAIN)

    alert({ header: 'Succesful', body: 'Profile is updated' })
    event.target.reset()
    setLoading(false)
  }

  return (
    <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div className="bg-white py-4 sm:py-8 px-4 shadow rounded-lg sm:px-10">
        <div className="text-2xl font-semibold mb-4">Update Profile</div>
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
                defaultValue={user?.email}
                disabled
                className="appearance-none block w-full px-3 py-2 border-2 bg-gray-100 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
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
                defaultValue={user?.displayName}
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
                defaultValue={user?.phone}
                required
                className="appearance-none block w-full px-3 py-2 border-2 border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-base sm:text-sm"
              />
            </div>
          </div>

          <div className="flex flex-row items-center gap-2">
            <input
              type="checkbox"
              id="contact_with_phone"
              name="contact_with_phone"
              value="true"
              defaultChecked={user?.contactWithPhone}
              className="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded"
            />
            <label
              htmlFor="contact_with_phone"
              className="text-md sm:text-sm text-gray-700"
            >
              Let other users contact me using my phone number
            </label>
          </div>

          {!loading && (
            <div>
              <button
                type="submit"
                className="w-full flex items-center gap-2 justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base sm:text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Update
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
