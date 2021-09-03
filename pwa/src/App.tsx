import React from 'react'
import { Route, Switch, Link, Redirect } from 'react-router-dom'
import styled from 'styled-components'

import CONSTANTS from './constants'
import MainView from './views/MainView'
import CreateRideView from './views/CreateRideView'
import RideDetailsView from './views/RideDetailsView'
import LoginView from './views/LoginView'
import RegisterView from './views/RegisterView'
import ProfileView from './views/ProfileView'
import Navbar from './components/Navbar'
import LogoutView from './views/LogoutView'
import useAuth from './hooks/useAuth'

export default function App() {
  const { user } = useAuth()

  return (
    <StyledContainer>
      <Link to="/">
        <h1>Country Roads</h1>
      </Link>
      <Navbar />

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
        <Route path={CONSTANTS.ROUTES.LOGOUT}>
          <LogoutView />
        </Route>
        <Route path={CONSTANTS.ROUTES.ME}>
          {!user && <Redirect to={CONSTANTS.ROUTES.RIDES.MAIN} />}
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
