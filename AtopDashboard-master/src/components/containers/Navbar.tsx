// @flow
import React, { useEffect } from "react";
import { Typography } from "antd";
import { useSelector, useDispatch } from "react-redux";
import AvatarDropdown from "components/purecomponents/AvatorDropdown";
import { clearAuth } from "store/auth/actions";
import { webVersion, useServerVersion } from "src/version";
import { RootState } from "src/store";
import { path } from "ramda";
import "./navbar.less";
const Navbar = () => {
  const dispatch = useDispatch();

  const doLogout = () => dispatch(clearAuth());

  const alreadyLogin = useSelector(path(["auth", "login"])) as boolean;
  const name = useSelector(path(["auth", "user", "name"])) as string;
  console.log("name", name);
  const serverVersion = useServerVersion();
  return (
    <header className="header">
      <div className="logo">
        <Typography.Title level={3}>Atop IOT server</Typography.Title>
      </div>
      <nav className="nav">
        <div className="navitem">
          <AvatarDropdown
            alreadyLogin={alreadyLogin}
            username={name}
            doLogout={doLogout}
          />
        </div>

        <div className="navitem">
          Web:{webVersion} | Server:{serverVersion}
        </div>
      </nav>
    </header>
  );
};

export default Navbar;
