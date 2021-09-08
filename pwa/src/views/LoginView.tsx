import React from 'react'
import LoginForm from '../components/LoginForm'
import OTPForm from '../components/OTPForm'
import RegisterForm from '../components/RegisterForm'

export default function LoginView() {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <LoginForm />
        <RegisterForm />
        <OTPForm />
      </div>
    </div>
  )
}
