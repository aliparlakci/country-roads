/* This example requires Tailwind CSS v2.0+ */
import React, { Fragment, useState, createContext, useContext } from 'react'
import cn from 'classnames'
import { Transition } from '@headlessui/react'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/outline'
import { XIcon } from '@heroicons/react/solid'

interface IModalData {
  header: string
  body: string
}

export interface IModalInterface {
  alert: (args: IModalData) => void
  error: (args: IModalData) => void
}

const modalContext = createContext<IModalInterface>({
  alert: (args: IModalData) => undefined,
  error: (args: IModalData) => undefined,
})

export function ModalProvider(props: any) {
  const [show, setShow] = useState(false)
  const [type, setType] = useState('')
  const [data, setData] = useState<IModalData>()

  const notification = (type: string) => (data: IModalData) => {
    setType(type)
    setData(data)
    setShow(true)
    setTimeout(() => setShow(false), 3000)
  }

  return (
    <modalContext.Provider
      value={{ alert: notification('alert'), error: notification('error') }}
    >
      {props.children}
      <div
        aria-live="assertive"
        className="z-50 fixed inset-0 flex items-end px-4 py-6 pointer-events-none sm:p-6 sm:items-start"
      >
        <div className="w-full flex flex-col items-center space-y-4 sm:items-end">
          {/* Notification panel, dynamically insert this into the live region when it needs to be displayed */}
          <Transition
            show={show}
            as={Fragment}
            enter="transform ease-out duration-300 transition"
            enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
            enterTo="translate-y-0 opacity-100 sm:translate-x-0"
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden">
              <div className="p-4">
                <div className="flex items-start">
                  <div className="flex-shrink-0">
                    {type === 'alert' && (
                      <CheckCircleIcon
                        className={cn('h-6 w-6 text-green-400')}
                        aria-hidden="true"
                      />
                    )}
                    {type === 'error' && (
                      <XCircleIcon
                        className={cn('h-6 w-6 text-red-400')}
                        aria-hidden="true"
                      />
                    )}

                  </div>
                  <div className="ml-3 w-0 flex-1 pt-0.5">
                    <p className="text-sm font-medium text-gray-900">
                      {data?.header}
                    </p>
                    <p className="mt-1 text-sm text-gray-500">{data?.body}</p>
                  </div>
                  <div className="ml-4 flex-shrink-0 flex">
                    <button
                      className="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                      onClick={() => {
                        setShow(false)
                      }}
                    >
                      <span className="sr-only">Close</span>
                      <XIcon className="h-5 w-5" aria-hidden="true" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </modalContext.Provider>
  )
}

const useModal = () => useContext(modalContext)

export default useModal
