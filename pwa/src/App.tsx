import React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import CONSTANTS from './constants'
import useAuth from './hooks/useAuth'
import Layout from './components/Layout'

import MainView from './views/MainView'
import CreateRideView from './views/CreateRideView'
import RideDetailsView from './views/RideDetailsView'
import LoginView from './views/LoginView'
import LogoutView from './views/LogoutView'
import OTPView from './views/OTPView'
import SettingsView from './views/SettingsView'

export default function App() {
  const { user } = useAuth()

  return (
    <Switch>
      <Route path={CONSTANTS.ROUTES.SIGNIN}>
        <LoginView />
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
            {!user && <Redirect to="/" />}
            {user && <CreateRideView />}
          </Route>
          <Route path={CONSTANTS.ROUTES.RIDES.DETAILS}>
            <RideDetailsView />
          </Route>
          <Route path={CONSTANTS.ROUTES.RIDES.MAIN}>
            <MainView />
          </Route>
          <Route path={CONSTANTS.ROUTES.SETTINGS}>
            <SettingsView />
          </Route>
          <Route path={CONSTANTS.ROUTES.ME}>
            {!user && <Redirect to={CONSTANTS.ROUTES.RIDES.MAIN} />}
            <></>
          </Route>
          <Redirect to={CONSTANTS.ROUTES.RIDES.MAIN} />
        </Switch>
      </Layout>
    </Switch>
  )
}
