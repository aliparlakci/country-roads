import React from 'react'
import { Link } from 'react-router-dom'
import styled from 'styled-components'
import CONSTANTS from '../constants'
import useAuth from '../hooks/useAuth'

export default function Navbar() {
    const { user } = useAuth()

    return (
        <StyledNav>
            <StyledNavRow>
                {user && (
                    <StyledNavItem>
                        <Link to={CONSTANTS.ROUTES.RIDES.NEW}>Create Ride</Link>
                    </StyledNavItem>
                )}
                <StyledNavItem>
                    <Link to={CONSTANTS.ROUTES.RIDES.MAIN}>Home</Link>
                </StyledNavItem>
            </StyledNavRow>
            {!user && (
                <StyledNavRow>
                    <StyledNavItem>
                        <Link to={CONSTANTS.ROUTES.LOGIN}>Login</Link>
                    </StyledNavItem>
                    <StyledNavItem>
                        <Link to={CONSTANTS.ROUTES.REGISTER}>Register</Link>
                    </StyledNavItem>
                </StyledNavRow>
            )}
            {user && (
                <StyledNavRow>
                    <StyledNavItem>
                        <span>{user.email}</span>
                    </StyledNavItem>
                    <StyledNavItem>
                        <Link to={CONSTANTS.ROUTES.LOGOUT}>Logout</Link>
                    </StyledNavItem>
                </StyledNavRow>
            )}
        </StyledNav>
    )
}

const StyledNavRow = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    width: 20rem;
`

const StyledNav = styled.nav`
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 1rem;
`

const StyledNavItem = styled.div`
    display: flex;
    justify-content: center;
    white-space: nowrap;
    text-decoration: none;
`
