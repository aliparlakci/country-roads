import React from 'react'
import { Route, Switch, Link } from 'react-router-dom'
import styled from 'styled-components'

import CONSTANTS from './constants'
import MainView from './views/MainView'
import CreateRideView from './views/CreateRideView'
import RideDetailsView from './views/RideDetailsView'
import LoginView from './views/LoginView'
import RegisterView from './views/RegisterView'
import ProfileView from './views/ProfileView'

export default function App() {
    return (
        <StyledContainer>
            <h1>Country Roads</h1>
            <StyledNav>
                <StyledNavItem>
                    <Link to={CONSTANTS.ROUTES.RIDES.NEW}>Create Ride</Link>
                </StyledNavItem>
                <StyledNavItem>
                    <Link to={CONSTANTS.ROUTES.RIDES.MAIN}>Home</Link>
                </StyledNavItem>
                <StyledNavItem>
                    <Link to={CONSTANTS.ROUTES.LOGIN}>Login</Link>
                </StyledNavItem>
                <StyledNavItem>
                    <Link to={CONSTANTS.ROUTES.REGISTER}>Register</Link>
                </StyledNavItem>
            </StyledNav>

            <Switch>
                <Route exact path={CONSTANTS.ROUTES.RIDES.NEW}>
                    <CreateRideView />
                </Route>
                <Route path={CONSTANTS.ROUTES.RIDES.DETAILS}>
                    <RideDetailsView />
                </Route>
                <Route path={CONSTANTS.ROUTES.RIDES.MAIN}>
                    <MainView />
                </Route>
                <Route path={CONSTANTS.ROUTES.LOGIN}>
                    <LoginView />
                </Route>
                <Route path={CONSTANTS.ROUTES.REGISTER}>
                    <RegisterView />
                </Route>
                <Route path={CONSTANTS.ROUTES.ME}>
                    <ProfileView />
                </Route>
            </Switch>
        </StyledContainer>
    )
}

const StyledContainer = styled.div`
    display: flex;
    justify-items: center;
    align-items: center;
    flex-direction: column;
    width: 100%;
    gap: 1rem;
`

const StyledNav = styled.nav`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    width: 20rem;
`

const StyledNavItem = styled.div`
    display: flex;
    justify-content: center;
    white-space: nowrap;
    text-decoration: none;
`
