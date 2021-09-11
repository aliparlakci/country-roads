import OTPForm from '../components/OTPForm'
import { Redirect, useLocation } from 'react-router-dom'
import CONSTANTS from '../constants'

export default function OTPView() {
  const { search } = useLocation()
  const params = new URLSearchParams(search)
  const email = params.get('email')

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        {!email && <Redirect to={CONSTANTS.ROUTES.SIGNIN} /> }
        <OTPForm email={email || ""} redirect={params.get('redirect') || ""} />
      </div>
    </div>
  )
}
