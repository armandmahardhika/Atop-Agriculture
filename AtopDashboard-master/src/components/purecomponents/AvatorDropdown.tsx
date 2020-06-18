import React from "react";
import { UserOutlined } from "@ant-design/icons";
import { Dropdown, Menu, Button } from "antd";
import { Link } from "react-router-dom";

type Props = {
  alreadyLogin: boolean;
  username: string;
  doLogout: Function;
  changePassword: Function;
};
const avatarDropdown = (props: Props) => {
  const { username, alreadyLogin, doLogout, changePassword } = props;

  const handleMenuClick = (k: { key: string }) => {
    const { key } = k;
    if (key === "logout") {
      console.log("key =", key);
      doLogout();
    }
  };

  const menu = (
    <Menu onClick={handleMenuClick}>
      <Menu.Item key="changepassword">
        <Link to="/resetpassword">Change Password</Link>
      </Menu.Item>
      <Menu.Item key="logout">Logout</Menu.Item>
    </Menu>
  );

  return (
    <div>
      {alreadyLogin ? (
        <Dropdown trigger={["click"]} overlay={menu}>
          <div>
            <Button icon={<UserOutlined />}>{username}</Button>
          </div>
        </Dropdown>
      ) : (
        <Link to="/login">Login</Link>
      )}
    </div>
  );
};

avatarDropdown.defaultProps = {
  doLogout: () => {},
  changePassword: () => {},
} as Partial<Props>;

export default avatarDropdown;
