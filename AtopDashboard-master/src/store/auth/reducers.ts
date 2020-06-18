import { SET_AUTH, CLEAR_AUTH } from "./actions";
import { AuthState, AuthActions } from "./types";

const initialState: AuthState = { login: false };

export function auth(
  state: AuthState = initialState,
  action: AuthActions
): AuthState {
  const { type } = action;
  switch (type) {
    case SET_AUTH:
      if ("payload" in action) {
        return { ...action.payload };
      }
      return state;
    case CLEAR_AUTH:
      return { login: false };
    default:
      return state;
  }
}
