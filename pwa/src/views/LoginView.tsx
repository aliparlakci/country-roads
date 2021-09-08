import React from 'react'
import { useLocation } from 'react-router-dom'
import LoginForm from '../components/LoginForm'

export default function LoginView() {
  const { search } = useLocation()
  const params = new URLSearchParams(search)
  const email = params.get('email')

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <LoginForm email={email || ""} />
      </div>
    </div>
  )
}
