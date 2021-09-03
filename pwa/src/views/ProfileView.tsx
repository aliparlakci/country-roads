import React from 'react'
import useAuth from '../hooks/useAuth'

export default function ProfileView() {
    const { user } = useAuth()

    return (
        <>
            <h3>ProfileView</h3>
            <span>{user?.displayName}</span>
            <span>{user?.email}</span>
            <span>{user?.phone}</span>
        </>
    )
}
