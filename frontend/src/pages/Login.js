import React, { useState, useMemo } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { ModalForgotPassword } from "../components/modals/ForgotPassword";
import {Login} from "../requests/Users";

function useQuery() {
  const { search } = useLocation();

  return useMemo(() => new URLSearchParams(search), [search]);
}

const SignInForm = () => {

  const [forgotPasswordShow, setForgotPasswordShow] = useState(false);
  const toggleForgotPassword = () => {
      setForgotPasswordShow(!forgotPasswordShow);
  };

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const [redirect, setRedirect] = useState(false);

  const [message, setMessage] = useState('');
  const msg = useQuery().get("message");

  if (window.localStorage.getItem("token") !== null) return <Navigate to="/"/>;

  const submit = async (e) => {
    e.preventDefault();

    try {
      if (email.trim() === "") {
          setMessage("Email cannot be blank");
          return;
      }
      if (password.trim() === "") {
          setMessage("Password cannot be blank");
          return;
      }
      let result = await Login(email, password);
      if (result.status === 200) {
        window.localStorage.setItem("token", result.data.accessToken);
        setRedirect(true);
      }

    } catch(err) {
      if (err.response !== undefined) setMessage(err.response.data.err);
      else setMessage(err.message);
    }
  }

  if (redirect) {
    return <Navigate to="/" />;
  }

  return (
    <div className="container pt-4 mt-4">
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
          <form onSubmit={submit}>
            <div className="display-3 mb-4 text-center text-primary"> Log In</div>
            <div id="message" className="mb-3 text-center text-danger">{message!==""?message:msg}</div>
            <div className="mb-3 text-center">
              <label htmlFor="email" className="form-label">Email address</label>
              <input type="email" className="form-control" id="email" onChange={e => setEmail(e.target.value)}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="password" className="form-label">Password</label>
              <input type="password" className="form-control" id="password" onChange={e => setPassword(e.target.value)}/>
            </div>
            <div className="text-center">
              <button type="submit" className="btn btn-primary text-center">Submit</button>
            </div>
            <div className="display-5 mt-4 text-center text-primary"><a href="/register">or Register</a></div>
            <div className="display-5 mt-4 text-center text-primary"><a href="#" onClick={toggleForgotPassword}>Forgot your password?</a></div>
          </form>
        </div>
      </div>
      <ModalForgotPassword show={forgotPasswordShow} toggle={toggleForgotPassword} />
    </div>
   )
}

export default SignInForm;