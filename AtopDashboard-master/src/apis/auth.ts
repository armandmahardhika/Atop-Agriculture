import React, { useState } from "react";
import axios from "axios";
import { JSONObject } from "./types";
type LoginResult = {
  data: any;
  loading: boolean;
  err: Error | undefined;
};

type LoginHookResult = [
  LoginResult,
  (form: { name: string; password: string }) => Promise<void>
];
// login hook send /api/token to IOTServer
export const useLogin = (): LoginHookResult => {
  const [loading, SetLoading] = useState(false);
  const [data, SetData] = useState<any>(undefined);
  const [err, SetErr] = useState<Error | undefined>(undefined);
  async function fetch(form: { name: string; password: string }) {
    try {
      SetLoading(true);
      const res = await axios.post("/api/token", form);
      SetData(res.data);
    } catch (err) {
      SetErr(err);
    } finally {
      SetLoading(false);
    }
  }

  return [{ data, loading, err }, fetch];
};
