import React from "react";
import { Route, Switch } from "react-router-dom";
import Testing from "components/pages/Testing";
import LoginPage from "components/pages/LoginPage";
const Routes = () => {
  return (
    <>
      <Switch>
        <Route path="/testing">
          <Testing />
        </Route>
        <Route path="/login" component={LoginPage} />
      </Switch>
    </>
  );
};

export default Routes;
