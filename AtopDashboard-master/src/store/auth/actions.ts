import { Token, User, AuthActions } from "./types";

import axios from "axios";
export const SET_AUTH = "SET_AUTH";
export const SET_USER = "SET_USER";
export const CLEAR_AUTH = "CLEAR_AUTH";

type SetAuthArg = {
  token: Token;
  user: User;
};
// Action factory
export const setAuth = (data: SetAuthArg): AuthActions => ({
  type: SET_AUTH,
  payload: { login: true, ...data },
});

export const clearAuth = (): AuthActions => ({ type: CLEAR_AUTH });
