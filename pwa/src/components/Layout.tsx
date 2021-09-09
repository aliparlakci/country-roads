import { Fragment, useState } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import { MenuIcon, XIcon } from '@heroicons/react/outline'
import Sidenav from '../components/Sidenav'
import CONSTANTS from '../constants'
import { Link } from 'react-router-dom'
import useAuth from '../hooks/useAuth'

export default function Layout({ children }: any) {
  const { user } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  return (
    <div
      className="flex bg-gray-100"
      style={{ height: 'min-content', minHeight: '100vh' }}
    >
      <div className="hidden lg:flex md:flex-shrink-0 w-64 h-screen bg-indigo-700">
        <Sidenav dark />
      </div>
      <Transition.Root show={sidebarOpen} as={Fragment}>
        <Dialog
          as="div"
          className="fixed inset-0 flex z-40 lg:hidden"
          onClose={setSidebarOpen}
        >
          <Transition.Child
            as={Fragment}
            enter="transition-opacity ease-linear duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="transition-opacity ease-linear duration-300"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Dialog.Overlay className="fixed inset-0 bg-gray-600 bg-opacity-75" />
          </Transition.Child>
          <Transition.Child
            as={Fragment}
            enter="transition ease-in-out duration-300 transform"
            enterFrom="-translate-x-full"
            enterTo="translate-x-0"
            leave="transition ease-in-out duration-300 transform"
            leaveFrom="translate-x-0"
            leaveTo="-translate-x-full"
          >
            <div className="relative flex-1 flex flex-col max-w-xs w-full">
              <Transition.Child
                as={Fragment}
                enter="ease-in-out duration-300"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="ease-in-out duration-300"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <div className="absolute top-0 right-0 -mr-12 pt-2">
                  <button
                    className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                    onClick={() => setSidebarOpen(false)}
                  >
                    <span className="sr-only">Close sidebar</span>
                    <XIcon className="h-6 w-6 text-white" aria-hidden="true" />
                  </button>
                </div>
              </Transition.Child>
              <div className="bg-indigo-700 h-full">
                <Sidenav dark={true} />
              </div>
            </div>
          </Transition.Child>
          <div className="flex-shrink-0 w-14" aria-hidden="true">
            {/* Dummy element to force sidebar to shrink to fit close icon */}
          </div>
        </Dialog>
      </Transition.Root>

      {/* Static sidebar for desktop */}
      <div className="flex flex-col w-screen flex-1">
        <div className="relative z-10 flex-shrink-0 flex lg:hidden items-center h-16 bg-white shadow">
          <button
            type="button"
            className="px-4 border-r border-gray-200 text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 lg:hidden"
            onClick={() => setSidebarOpen(true)}
          >
            <span className="sr-only">Open sidebar</span>
            <MenuIcon className="h-6 w-6" aria-hidden="true" />
          </button>

          <div className="flex-1 mr-8 2xl:mr-64 flex justify-end">
            <div className="flex items-center">
              {!user && (
                <Link
                  to={CONSTANTS.ROUTES.LOGIN}
                  className="bg-indigo-600 hover:bg-indigo-800 text-white group flex items-center px-4 py-2 text-base font-medium rounded-full shadow max-h-12 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 border border-transparent"
                >
                  Sign in
                </Link>
              )}
            </div>
          </div>
        </div>
        <div className="flex flex-row justify-center items-start h-full w-full">
          <main className="py-6 px-4 w-full h-full sm:px-6 md:px-8 max-w-4xl">
            {children}
          </main>
        </div>
      </div>
    </div>
  )
}
