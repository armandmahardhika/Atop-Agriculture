import React, { useState, useEffect } from "react";
import axios from "axios";

export const useServerVersion = () => {
  const [version, setVersion] = useState("");

  useEffect(() => {
    axios
      .get("api/version")
      .then((res) => setVersion(res.data.version))
      .catch((err) => setVersion("Unknown"));
  }, []);
  return version;
};

export const webVersion = require("../package.json").version;
