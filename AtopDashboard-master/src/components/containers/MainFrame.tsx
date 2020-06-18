import React, { useState, useEffect } from "react";
import { Layout } from "antd";
import Navbar from "components/containers/Navbar";
import Routes from "routers/Routes";
import "./MainFrame.less";
const mainFrame = () => {
  return (
    <div>
      <Navbar />
      <div className="content">
        <Routes />
      </div>
    </div>
  );
};

export default mainFrame;
