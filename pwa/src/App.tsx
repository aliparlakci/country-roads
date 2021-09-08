import React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import CONSTANTS from './constants'
import MainView from './views/MainView'
import CreateRideView from './views/CreateRideView'
import RideDetailsView from './views/RideDetailsView'
import LoginView from './views/LoginView'
import RegisterView from './views/RegisterView'
import ProfileView from './views/ProfileView'
import LogoutView from './views/LogoutView'
import useAuth from './hooks/useAuth'
import Layout from './components/Layout'
import OTPView from './views/OTPView'

export default function App() {
  const { user } = useAuth()

  return (
    <Switch>
      <Route path={CONSTANTS.ROUTES.LOGIN}>
        <LoginView />
      </Route>
      <Route path={CONSTANTS.ROUTES.REGISTER}>
        <RegisterView />
      </Route>
      <Route path={CONSTANTS.ROUTES.OTP}>
        <OTPView />
      </Route>
      <Route path={CONSTANTS.ROUTES.LOGOUT}>
        <LogoutView />
      </Route>

      <Layout>
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
          <Route path={CONSTANTS.ROUTES.ME}>
            {!user && <Redirect to={CONSTANTS.ROUTES.RIDES.MAIN} />}
            <ProfileView />
          </Route>
          <Redirect to={CONSTANTS.ROUTES.RIDES.MAIN} />
        </Switch>
      </Layout>
    </Switch>
  )
}
