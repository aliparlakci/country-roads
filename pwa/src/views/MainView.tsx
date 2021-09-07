import React from 'react'
import styled from 'styled-components'

import RideList from '../components/RideList'
import useQuery from '../hooks/useQuery'
import { IRideQuery } from '../hooks/useRides'
import RideFilter from '../components/RideFilter'
import { Disclosure } from '@headlessui/react'
import { ArrowDownIcon, ArrowUpIcon } from '@heroicons/react/outline'

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
    <StyledContainer>
      <Disclosure>
        <Disclosure.Button
          as="button"
          className="flex flex-row items-center gap-1 text-left text-2xl font-semibold text-gray-900"
        >
          {({ open }) => (
            <div className="flex items-center gap-1 ml-2 text-blue-900">  
              Filters
              {!open && <ArrowDownIcon className="h-4" />}
              {open && <ArrowUpIcon className="h-4" />}
            </div>
          )}
        </Disclosure.Button>
        <Disclosure.Panel>
          <RideFilter />
        </Disclosure.Panel>
      </Disclosure>
      <RideList {...query} />
    </StyledContainer>
  )
}

const StyledContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;

  height: 100%;
  width: 100%;
`

const ColumnView = styled.div`
  display: flex;
  flex-direction: row;
  gap: 2rem;
`
