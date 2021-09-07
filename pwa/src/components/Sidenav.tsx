import React from 'react'
import {
  AdjustmentsIcon,
  KeyIcon,
  LoginIcon,
  LogoutIcon,
  MailIcon,
  NewspaperIcon,
  UserIcon,
} from '@heroicons/react/outline'
import { Link } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGithub } from '@fortawesome/free-brands-svg-icons'
import useAuth from '../hooks/useAuth'
import CONSTANTS from '../constants'

function classNames(...classes: any[]) {
  return classes.filter(Boolean).join(' ')
}

export default function Sidenav() {
  return (
    <>
      <div className="flex flex-col space-y-4 w-full h-full pt-5 bg-indigo-700">
        <div className="flex-1 h-0 overflow-y-auto">
          <Navigation />
        </div>
        <Credits />
        <UserInfo />
      </div>
    </>
  )
}

interface IMenuItemProps {
  label: string
  href: string
  current: boolean
  Icon: (props: React.ComponentProps<'svg'>) => JSX.Element
}

function Navigation() {
  const { user } = useAuth()

  return (
    <nav className="px-2 space-y-1">
      <MenuItem
        href={CONSTANTS.ROUTES.RIDES.MAIN}
        current={false}
        Icon={NewspaperIcon}
        label="Feed"
      />
      {user && (
        <>
          <MenuItem
            href={CONSTANTS.ROUTES.ME}
            current={false}
            Icon={UserIcon}
            label="Profile"
          />
          <MenuItem
            href="?"
            current={false}
            Icon={AdjustmentsIcon}
            label="Settings"
          />
        </>
      )}
    </nav>
  )
}

function MenuItem({ label, Icon, href, current }: IMenuItemProps) {
  return (
    <Link
      key={label}
      to={href}
      className={classNames(
        current
          ? 'bg-indigo-800 text-white'
          : 'text-indigo-100 hover:bg-indigo-600',
        'group flex items-center px-2 py-2 text-base font-medium rounded-md',
      )}
    >
      <Icon
        className="mr-4 flex-shrink-0 h-6 w-6 text-indigo-300"
        aria-hidden="true"
      />
      {label}
    </Link>
  )
}

function Credits() {
  return (
    <div className="mx-5 space-y-6">
      <div className="text-sm text-left text-indigo-300 select-none">
        This project is made by <b>Ali Parlakçı</b>
      </div>
      <div className="gap-2 justify-center flex flex-col">
        <div className="flex items-center">
          <FontAwesomeIcon
            icon={faGithub}
            className="text-indigo-300 text-xl"
          />
          <span className="text-sm pl-2 text-indigo-300">
            github.com/aliparlakci
          </span>
        </div>
        <div className="flex items-start">
          <MailIcon className="inline-block text-indigo-300 h-5" />
          <span className="text-sm pl-2 text-indigo-300">
            aliparlakci@sabanciuniv.edu
          </span>
        </div>
      </div>
    </div>
  )
}

function UserInfo() {
  const { user } = useAuth()

  if (user)
    return (
      <div className="flex flex-col gap-2 items-center border-t border-indigo-800 p-4">
        <div className="w-full">
          <div className="flex items-center">
            <div>
              <p className="text-sm font-medium text-white">
                {user.displayName}
              </p>
              <p className="text-xs font-medium text-indigo-200">
                {user.email}
              </p>
            </div>
          </div>
        </div>
        <Link to={CONSTANTS.ROUTES.LOGOUT} className="w-full flex justify-center items-center text-sm h-8 rounded-full border-2 border-white hover:border-red-500 text-white hover:text-red-500 transition">
          Sign out
        </Link>
      </div>
    )

  return <div className="h-4"></div>
}
