import React, { createContext, useContext } from 'react'
import useSWR, { mutate } from 'swr'
import CONSTANTS from '../constants'
import { IUser } from '../types/user'

export interface IAuthContext {
    user?: IUser
    logout: CallableFunction
}

const authContext = createContext<IAuthContext>({ logout: () => null })

export const logout = async () => {
    await fetch(CONSTANTS.API.AUTH.LOGOUT, { method: 'POST' })
}

export function AuthProvider(props: any) {
    const { data } = useSWR<{ user?: IUser, error: string }>(CONSTANTS.API.AUTH.USER)

    return (
        <authContext.Provider value={{ user: data?.user, logout }}>
            {props.children}
        </authContext.Provider>
    )
}

const useAuth = () => useContext(authContext)

export default useAuth
