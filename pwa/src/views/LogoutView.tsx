import React, { useState, useEffect } from 'react'
import { Redirect } from 'react-router-dom'
import { mutate } from 'swr'
import CONSTANTS from '../constants'
import useAuth from '../hooks/useAuth'

export default function LogoutView() {
  const { logout } = useAuth()
  const [redirect, setRedirect] = useState(false)

  useEffect(() => {
    const doLogout = async () => {
      await logout()
      setRedirect(true)
      mutate(CONSTANTS.API.AUTH.USER)
    }

    doLogout()
  }, [logout])

  return <>{redirect && <Redirect to="/"></Redirect>}</>
}
