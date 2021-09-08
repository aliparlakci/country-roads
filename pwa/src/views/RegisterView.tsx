import React from 'react'
import { useLocation } from 'react-router-dom'
import RegisterForm from '../components/RegisterForm'

export default function RegisterView() {
  const { search } = useLocation()
  const params = new URLSearchParams(search)

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <RegisterForm email={params.get('email') || ''} />
      </div>
    </div>
  )
}
