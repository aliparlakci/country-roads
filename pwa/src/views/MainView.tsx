import React from 'react'
import { Disclosure } from '@headlessui/react'
import { ArrowDownIcon, ArrowUpIcon } from '@heroicons/react/outline'

import RideList from '../components/RideList'
import useQuery from '../hooks/useQuery'
import { IRideQuery } from '../hooks/useRides'
import RideFilter from '../components/RideFilter'

export default function MainView() {
  const params = useQuery()
  const query: IRideQuery = {
    type: params.get('type'),
    direction: params.get('direction'),
    destination: params.get('destination'),
    startDate: params.get('start_date'),
    endDate: params.get('end_date'),
  }

  return (
    <div className="flex flex-col gap-4 h-full w-full">
      <Disclosure>
        <Disclosure.Button
          as="button"
          className="flex flex-row items-center gap-1 text-left text-2xl font-semibold text-gray-900"
        >
          {({ open }) => (
            <div className="flex items-center gap-1 ml-2 text-black">  
              Filters
              {!open && <ArrowDownIcon className="animate-bounce h-4" />}
              {open && <ArrowUpIcon className="h-4" />}
            </div>
          )}
        </Disclosure.Button>
        <Disclosure.Panel>
          <RideFilter />
        </Disclosure.Panel>
      </Disclosure>
      <RideList {...query} />
    </div>
  )
}