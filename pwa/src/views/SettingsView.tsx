import React from 'react'
import UpdateProfile from '../components/UpdateProfile'

export default function SettingsView() {
  return (
    <>
      <div className="flex flex-col gap-4 w-full h-full">
        <div className="flex flex-row items-center gap-1 text-left text-4xl font-semibold text-gray-800 ml-2">
          Settings
        </div>
        <UpdateProfile />
      </div>
    </>
  )
}
