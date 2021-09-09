import React from 'react'
import {
  AdjustmentsIcon,
  LoginIcon,
  MailIcon,
  NewspaperIcon,
  PlusIcon,
  UserIcon,
} from '@heroicons/react/outline'
import { Link } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGithub } from '@fortawesome/free-brands-svg-icons'
import cn from 'classnames'
import useAuth from '../hooks/useAuth'
import CONSTANTS from '../constants'

export default function Sidenav({ dark }: { dark?: boolean }) {
  return (
    <>
      <div className="flex flex-col space-y-4 w-full h-full pt-5">
        <div className="flex-1 h-0 overflow-y-auto">
          <Navigation dark={dark} />
        </div>
        <Credits dark={dark} />
        <UserInfo dark={dark} />
      </div>
    </>
  )
}

interface IMenuItemProps {
  dark?: boolean
  label: string
  href: string
  current: boolean
  Icon: (props: React.ComponentProps<'svg'>) => JSX.Element
}

function Navigation({ dark }: { dark?: boolean }) {
  const { user } = useAuth()

  return (
    <nav className="px-2 space-y-1">
      {!user && (
        <MenuItem
          href={CONSTANTS.ROUTES.LOGIN}
          current={false}
          Icon={LoginIcon}
          label="Sign in"
          dark={true}
        />
      )}
      {user && (
        <MenuItem
          href={CONSTANTS.ROUTES.RIDES.NEW}
          current={false}
          Icon={PlusIcon}
          label="New post"
          dark={true}
        />
      )}
      <MenuItem
        href={CONSTANTS.ROUTES.RIDES.MAIN}
        current={false}
        Icon={NewspaperIcon}
        label="Feed"
        dark={dark}
      />
      {user && (
        <>
          <MenuItem
            href={CONSTANTS.ROUTES.ME}
            current={false}
            Icon={UserIcon}
            label="Profile"
            dark={dark}
          />
        </>
      )}
    </nav>
  )
}

function MenuItem({ label, Icon, href, current, dark }: IMenuItemProps) {
  return (
    <Link
      key={label}
      to={href}
      className={cn(
        'group flex items-center px-2 py-2 text-base font-medium rounded-md hover:bg-indigo-600',
        { 'text-indigo-200 hover:text-white': dark },
        { 'text-black hover:text-white': !dark },
      )}
    >
      <Icon
        className={cn(
          'mr-4 flex-shrink-0 h-6 w-6 group-hover:text-white',
          { 'text-indigo-200': dark },
          { 'text-black': !dark },
        )}
        aria-hidden="true"
      />
      {label}
    </Link>
  )
}

function Credits({ dark }: { dark?: boolean }) {
  return (
    <div className="mx-5 space-y-6">
      <div
        className={cn(
          { 'text-indigo-200': dark },
          { 'text-gray-500': !dark },
          'text-sm text-left select-none',
        )}
      >
        {/* This project is developed by <b>Ali Parlakçı</b> */}
      </div>
      <div className="gap-2 justify-center flex flex-col">
        <div className="flex items-center">
          <FontAwesomeIcon
            icon={faGithub}
            className={cn(
              'text-xl',
              { 'text-indigo-200': dark },
              { 'text-gray-500': !dark },
            )}
          />
          <a
            href="https://github.com/aliparlakci"
            className={cn(
              { 'text-indigo-200': dark },
              { 'text-gray-500': !dark },
              'text-sm pl-2',
            )}
          >
            github.com/aliparlakci
          </a>
        </div>
        <div className="flex items-start">
          <MailIcon
            className={cn(
              'inline-block h-5',
              { 'text-indigo-200': dark },
              { 'text-gray-500': !dark },
            )}
          />
          <a
            href="mailto:aliparlakci@sabanciuniv.edu"
            className={cn(
              'text-sm pl-2',
              { 'text-indigo-200': dark },
              { 'text-gray-500': !dark },
            )}
          >
            aliparlakci@sabanciuniv.edu
          </a>
        </div>
      </div>
    </div>
  )
}

function UserInfo({ dark }: { dark?: boolean }) {
  const { user } = useAuth()

  if (user)
    return (
      <div
        className={cn(
          'flex flex-col gap-2 items-center border-t p-4',
          { 'border-indigo-800': dark },
          { 'border-gray-300': !dark },
        )}
      >
        <div className="w-full">
          <div className="flex items-center">
            <div>
              <p
                className={cn(
                  'text-sm font-medium',
                  { 'text-white': dark },
                  { 'text-black': !dark },
                )}
              >
                {user.displayName}
              </p>
              <p
                className={cn(
                  'text-xs font-medium',
                  { 'text-indigo-200': dark },
                  { 'text-black': !dark },
                )}
              >
                {user.email}
              </p>
            </div>
          </div>
        </div>
        <Link
          to={CONSTANTS.ROUTES.LOGOUT}
          className={cn(
            'transition w-full flex justify-center items-center text-sm h-8 rounded-full border-2  hover:border-red-500 font-medium hover:text-red-500',
            { 'text-white border-white': dark },
            { 'text-black border-black': !dark },
          )}
        >
          Sign out
        </Link>
      </div>
    )

  return <div className="h-4"></div>
}
