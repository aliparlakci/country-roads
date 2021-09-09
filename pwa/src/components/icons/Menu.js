import * as React from 'react'

function SvgMenu(props) {
  return (
    <svg width="1em" height="1em" viewBox="0 0 64 64" fill="none" {...props}>
      <path
        d="M8 32h48M8 16h48M8 48h48"
        stroke="currentColor"
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default SvgMenu
