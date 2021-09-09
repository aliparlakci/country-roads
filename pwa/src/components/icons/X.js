import * as React from 'react'

function SvgX(props) {
  return (
    <svg width="1em" height="1em" viewBox="0 0 64 64" fill="none" {...props}>
      <path
        d="M48 16L16 48M16 16l32 32"
        stroke="currentColor"
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default SvgX
