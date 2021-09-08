import React, { useEffect } from 'react'
import { faCircleNotch } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { useHistory } from 'react-router-dom'
import { mutate } from 'swr'
import CONSTANTS from '../constants'
import useAuth from '../hooks/useAuth'

export default function LogoutView() {
  const history = useHistory()
  const { logout } = useAuth()

  useEffect(() => {
    const doLogout = async () => {
      await logout()
      await mutate(CONSTANTS.API.AUTH.USER)
      history.push('/')
    }

    doLogout()
  }, [history, logout])

  return (
    <div className="flex justify-center items-center w-screen h-screen bg-gray-100">
      <FontAwesomeIcon
        className="animate-spin text-2xl text-indigo-600"
        icon={faCircleNotch}
      />
    </div>
  )
}
