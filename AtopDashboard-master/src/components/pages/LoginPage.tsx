import React, { useEffect } from "react";
import { path } from "ramda";
import LoginForm from "components/forms/LoginForm";
import { setAuth, clearAuth } from "store/auth/actions";
import { message } from "antd";
import { useDispatch } from "react-redux";
import { useLogin } from "apis/auth";
import { isFine, errorMessage } from "apis/errors";
import "./LoginPage.less";
const loginPage = () => {
  const dispatch = useDispatch();
  const [{ data, loading }, fetchToken] = useLogin();

  useEffect(() => {
    if (!data) {
      return;
    }
    if (!isFine(data)) {
      message.error(errorMessage(data));
      return;
    }

    const auth = {
      token: path(["payload", "token"], data),
      user: path(["payload", "user"], data),
    };
    //const auth = { token: undefined, user: undefined };
    dispatch(setAuth(auth as any));
  }, [data]);
  const onSubmit = (form: { name: string; password: string }) => {
    fetchToken({ ...form });
  };
  return (
    <div>
      <LoginForm onSubmit={onSubmit} loading={loading} />
      <button onClick={() => dispatch(clearAuth())}>logout</button>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};
export default loginPage;
