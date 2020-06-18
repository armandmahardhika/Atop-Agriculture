import { Action } from "redux";

export interface Token {
  token: string;
  refreshToken: string;
}

export interface User {
  name: string;
  id: string;
}

export interface AuthState {
  login: boolean;
  token?: Token;
  user?: User;
}

// action types
export interface SetAuth {
  type: string;
  payload: AuthState;
}

export type AuthActions = SetAuth | Action<string>;
