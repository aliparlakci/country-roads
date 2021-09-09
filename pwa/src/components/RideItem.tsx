import React from 'react'
import cn from 'classnames'

import CONSTANTS from '../constants'
import useAuth from '../hooks/useAuth'
import IRide from '../types/ride'
import mutateWithQueries from '../utils/mutateWithQueries'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
  faCircleNotch,
  faPhone,
  faEnvelope,
} from '@fortawesome/free-solid-svg-icons'
import IContactInfo from '../types/contact'

export interface IRideItemProps {
  ride: IRide
}

export default function RideItem({ ride }: IRideItemProps) {
  const [deleteLoading, setDeleteLoading] = React.useState(false)
  const [contact, setContact] = React.useState<IContactInfo | null>(null)
  const { user } = useAuth()

  const getContact = async (id: string) => {
    let response
    try {
      response = await fetch(`${CONSTANTS.API.USERS.CONTACT}/${id}`)
    } catch (e) {
      console.log(e)
    }

    setContact(await response?.json())
  }

  const doDelete = async () => {
    setDeleteLoading(true)

    let response
    try {
      response = await fetch(CONSTANTS.API.RIDE(ride.id), { method: 'delete' })
    } catch (err) {
      console.error(err)
    }

    if (!response || response?.status !== 200) {
      setDeleteLoading(false)
      // TODO: Bring up an alert
      return
    }

    await mutateWithQueries(CONSTANTS.API.RIDES)
  }

  return (
    <div>
      <div
        className={cn(
          'rounded-t-md bg-white border-l-4',
          { 'border-yellow-300': ride.type === 'taxi' },
          { 'border-green-400': ride.type === 'offer' },
          { 'border-blue-400': ride.type === 'request' },
        )}
      >
        <div className="border-b border-gray-200 text-2xl sm:text-xl font-light px-4 py-1">
          {(() => {
            switch (ride.type) {
              case 'offer':
                return 'Offers a lift'
              case 'request':
                return 'Needs a lift'
              case 'taxi':
                return 'Shares a taxi'
              default:
                return ''
            }
          })()}
        </div>
      </div>
      <div
        className={cn(
          'flex flex-col gap-2 bg-white border-l-4 rounded-b-md shadow-md p-3',
          { 'border-yellow-300': ride.type === 'taxi' },
          { 'border-green-400': ride.type === 'offer' },
          { 'border-blue-400': ride.type === 'request' },
        )}
      >
        <div
          className="grid gap-4 items-center ml-2"
          style={{ gridTemplateColumns: 'auto 1fr', rowGap: '.5rem' }}
        >
          <span className="text-base sm:text-sm font-semibold text-right">
            From
          </span>
          <span className="text-base sm:text-sm font-light">
            {ride.from.display}
          </span>
          <span className="text-base sm:text-sm font-semibold text-right">
            To
          </span>
          <span className="text-base sm:text-sm font-light">
            {ride.to.display}
          </span>

          <span className="text-base sm:text-sm font-semibold text-right">
            When
          </span>
          <span className="text-base sm:text-sm font-light">
            {(() => {
              const date = new Date(ride.date * 1000).toDateString()
              return date.substring(0, date.length - 5)
            })()}
          </span>
          <span className="text-base sm:text-sm font-semibold text-right">
            Author
          </span>
          <span className="text-base sm:text-sm font-light">
            {ride.owner.displayName}
          </span>
        </div>
        <div className="flex flex-row justify-between items-center">
          {user && (
            <div className="flex flex-row gap-1 justify-between">
              {!contact && (
                <button
                  className={cn(
                    'flex items-center justify-center text-base sm:text-sm rounded-full border-2 border-transparent px-2 mb-1 transition hover:text-white w-min m-0',
                    {
                      'text-yellow-500 hover:bg-yellow-600':
                        ride.type === 'taxi',
                    },
                    {
                      'text-green-600 hover:bg-green-600':
                        ride.type === 'offer',
                    },
                    {
                      'text-blue-500 hover:bg-blue-600':
                        ride.type === 'request',
                    },
                  )}
                  onClick={() => getContact(ride.owner.id)}
                >
                  Contact
                </button>
              )}
              {contact !== null && (
                <div className="ml-2 flex flex-row gap-6 text-lg text-indigo-600">
                  <a href={`mailto:${contact.email}`}>
                    <FontAwesomeIcon
                      icon={faEnvelope}
                      className={cn(
                        'transition',
                        {
                          'text-yellow-400 hover:text-yellow-600':
                            ride.type === 'taxi',
                        },
                        {
                          'text-green-400 hover:text-green-600':
                            ride.type === 'offer',
                        },
                        {
                          'text-blue-400 hover:text-blue-600':
                            ride.type === 'request',
                        },
                      )}
                    />
                  </a>
                  {contact.phone && (
                    <a href={`tel:${contact.phone}`}>
                      <FontAwesomeIcon
                        icon={faPhone}
                        className={cn(
                          'transition',
                          {
                            'text-yellow-400 hover:text-yellow-600':
                              ride.type === 'taxi',
                          },
                          {
                            'text-green-400 hover:text-green-600':
                              ride.type === 'offer',
                          },
                          {
                            'text-blue-400 hover:text-blue-600':
                              ride.type === 'request',
                          },
                        )}
                      />
                    </a>
                  )}
                </div>
              )}
            </div>
          )}
          {!user && (
            <div className="text-xs text-gray-400 font-extralight italic text-right w-full">
              Sign in to contact
            </div>
          )}
          {user?.id === ride.owner.id && (
            <>
              {!deleteLoading && (
                <button
                  className="flex items-center justify-center text-base sm:text-sm rounded-full border-2 border-transparent px-2 mb-1 transition hover:bg-red-600 text-red-600 hover:text-white w-min m-0"
                  onClick={doDelete}
                >
                  Delete
                </button>
              )}
              {deleteLoading && (
                <FontAwesomeIcon
                  className="animate-spin text-lg text-red-400"
                  icon={faCircleNotch}
                />
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
