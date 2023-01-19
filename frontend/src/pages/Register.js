import React, { useState } from "react";
import { Navigate } from "react-router-dom";
import {Register} from "../requests/Users";

function RegisterForm() {
  const [username, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rpassword, setRPassword] = useState("");
  const [message, setMessage] = useState("");

  if (window.localStorage.getItem("token") !== null) return <Navigate to="/"/>

  const submit = async (e) => {
    e.preventDefault();

    let result;
    try {
      if (email.trim() === "") {
          setMessage("Email can't be blank");
          return;
      }
      if (username.trim() === "") {
          setMessage("Username can't be blank");
          return;
      }
      if (password.trim() === "") {
          setMessage("Password can't be blank");
          return;
      }
      if (password !== rpassword) {
          setMessage("Passwords don't match");
          return;
      }
      result = await Register(email, username, password, rpassword);
    } catch(err) {
      setMessage(err.response.data.err);
    }
    if (result.status === 201) {
      setMessage("Verification email sent");
    }

  }

  return (
    <div className="container mt-4 pt-4">
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
          <form onSubmit={submit}>
            <div className="display-1 mb-4 text-center text-primary"> Register</div>
            <div className="mb-3 text-center text-danger">{message}</div>
            <div className="mb-3 text-center">
              <label htmlFor="username" className="form-label">Username</label>
              <input name="username" type="text" className="form-control" id="username" onChange={(e) => setName(e.target.value)}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="email" className="form-label">Email address</label>
              <input name="email" type="email" className="form-control" id="email" onChange={(e) => setEmail(e.target.value)}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="pass" className="form-label">Password</label>
              <input name="password" type="password" className="form-control" id="password" onChange={(e) => setPassword(e.target.value)}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="pass" className="form-label">Repeat Password</label>
              <input name="rpassword" type="password" className="form-control" id="rpassword" onChange={(e) => setRPassword(e.target.value)}/>
            </div>
            <div className="text-center">
              <button type="submit" className="btn btn-primary text-center">Submit</button>
            </div>
            <div className="display-5 mt-4 text-center text-primary"><a href="/login">or Log in</a></div>
          </form>
        </div>
      </div>
      </div>
    );
  }

export default RegisterForm;