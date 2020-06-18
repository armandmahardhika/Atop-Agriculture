import React from "react";
import { Router } from "react-router";
import { Provider } from "react-redux";
import history from "routers/history";
import { store } from "src/store";
import MainFrame from "components/containers/MainFrame";
import "./App";
const App = () => {
  return (
    <Provider store={store}>
      <Router history={history}>
        <MainFrame />
      </Router>
    </Provider>
  );
};

export default App;
