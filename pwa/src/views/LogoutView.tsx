import React, { useState, useEffect } from 'react'
import { useHistory, Redirect } from 'react-router-dom'
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

  return <>Loading</>
}
