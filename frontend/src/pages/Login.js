import React, { useState } from "react";
import { Navigate } from "react-router-dom";
import APICaller from "../Requests";

const SignInForm = (props) => {

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [redirect, setRedirect] = useState(false);
  const [message, setMessage] = useState('');

  const submit = async (e) => {
    e.preventDefault();
    let result;

    try {
      result = await APICaller.Login(email, password);
    } catch(err) {
      setMessage(err.message);
    }

    if (result.status !== 200) {
      setMessage(result.data.err);
      return;
    }
    APICaller.SetAccessToken(result.data.accessToken);
    props.setName("logged");
    setRedirect(true);
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
          <div id="message" className="mb-3 text-center text-danger">{message}</div>
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
        </form>
      </div>
    </div>
    </div>
   )
}

export default SignInForm;