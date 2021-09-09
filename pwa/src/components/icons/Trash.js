import * as React from 'react'

function SvgTrash(props) {
  return (
    <svg width="1em" height="1em" viewBox="0 0 25 25" fill="none" {...props}>
      <path
        d="M3.5 6.5h2m0 0h16m-16 0v14a2 2 0 002 2h10a2 2 0 002-2v-14h-14zm3 0v-2a2 2 0 012-2h4a2 2 0 012 2v2"
        stroke="currentColor"
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default SvgTrash
