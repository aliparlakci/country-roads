import * as React from 'react'

function SvgLink(props) {
  return (
    <svg width="1em" height="1em" viewBox="0 0 64 64" fill="none" {...props}>
      <path
        d="M28.782 33.61a8.047 8.047 0 005.868 3.21 8.038 8.038 0 006.264-2.34l4.828-4.832a8.058 8.058 0 00-.098-11.29 8.042 8.042 0 00-11.279-.097l-2.767 2.754"
        stroke="currentColor"
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M35.218 30.39a8.048 8.048 0 00-5.868-3.21 8.037 8.037 0 00-6.264 2.34l-4.828 4.832a8.058 8.058 0 00.098 11.29 8.042 8.042 0 0011.279.097l2.751-2.754"
        stroke="currentColor"
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default SvgLink
