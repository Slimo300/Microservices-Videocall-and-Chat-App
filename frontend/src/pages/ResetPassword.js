import React, { useState, useEffect } from "react";
import { Navigate, useParams} from "react-router-dom";
import {ResetForgottenPassword} from "../requests/Users";

function ResetPassword() {
  const [resetCode, setResetCode] = useState("");

  const [newPassword, setNewPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");

  const [message, setMessage] = useState("");
  const [redirect, setRedirect] = useState(false);
  const {code} = useParams();

  useEffect(() => {
    if (code !== "" ) {
      setResetCode(code);
    }
  }, [code]);

  if (window.localStorage.getItem("token") !== null) return <Navigate to="/"/>

  const submit = async (e) => {
    e.preventDefault();

    try {
      let result = await ResetForgottenPassword(resetCode, newPassword, repeatPassword);
      if (result.status === 200) {
        setMessage("Password changed");
      }
    } catch(err) {
      if (err.response !== undefined) setMessage(err.response.data.err);
      else setMessage(err.message);
    }
    
    setTimeout(function() {
      setMessage("");
      setRedirect(true);
  }, 2000);
  }

  if (redirect) {
    return <Navigate to="/login?message=Password+changed"/>;
  }

  return (
    <div className="container mt-4 pt-4">
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
          <form onSubmit={submit}>
            <div className="display-1 mb-4 text-center text-primary"> Reset your password</div>
            <div className="mb-3 text-center text-danger">{message}</div>
            <div className="mb-3 text-center">
              <label className="form-label">New Password:</label>
              <input name="password" type="password" className="form-control" id="verification-code" onChange={(e) => setNewPassword(e.target.value)}/>
            </div>
            <div className="mb-3 text-center">
              <label className="form-label">Repeat Password:</label>
              <input name="rpassword" type="password" className="form-control" id="verification-code" onChange={(e) => setRepeatPassword(e.target.value)}/>
            </div>
            <div className="text-center">
              <button type="submit" className="btn btn-primary text-center">Change Password</button>
            </div>
          </form>
        </div>
      </div>
      </div>
    );
  }

export default ResetPassword;