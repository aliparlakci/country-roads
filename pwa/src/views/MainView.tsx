import React from 'react'
import { Disclosure } from '@headlessui/react'
import cn from 'classnames'
import { Link } from 'react-router-dom'
import { PlusSmIcon } from '@heroicons/react/outline'

import RideList from '../components/RideList'
import useAuth from '../hooks/useAuth'
import useQuery from '../hooks/useQuery'
import { IRideQuery } from '../hooks/useRides'
import RideFilter from '../components/RideFilter'
import CONSTANTS from '../constants'

export default function MainView() {
  const params = useQuery()
  const { user } = useAuth()
  const query: IRideQuery = {
    type: params.get('type'),
    direction: params.get('direction'),
    destination: params.get('destination'),
    startDate: params.get('start_date'),
    endDate: params.get('end_date'),
  }

  return (
    <div className="flex flex-col gap-4 h-full w-full">
      <div className="flex flex-row items-center gap-1 text-left text-4xl font-semibold text-gray-800 ml-2">
        Feed
      </div>

      <Disclosure>
        <div className="w-full flex flex-row-reverse justify-between">
          <Disclosure.Button as="button">
            {({ open }) => (
              <div
                className={cn(
                  'transition px-3 py-1 flex justify-end items-center gap-1 mr-2 text-base sm:text-sm rounded-full border-2 border-transparent text-indigo-600',
                  { 'border-indigo-600': open },
                )}
              >
                Show filters
              </div>
            )}
          </Disclosure.Button>
          {user && (
            <Link
              to={CONSTANTS.ROUTES.RIDES.NEW}
              className="shadow bg-indigo-600 border-2 border-indigo-600 rounded-full px-3 py-1 ml-1 text-base sm:text-sm text-white flex justify-center items-center gap-1 hover:bg-indigo-700 hover:border-indigo-700 transition"
            >
              <PlusSmIcon className="h-6 sm:h-5" />
              New Post
            </Link>
          )}
        </div>
        <Disclosure.Panel>
          <RideFilter />
        </Disclosure.Panel>
      </Disclosure>
      <RideList {...query} />
    </div>
  )
}
