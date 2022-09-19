import React, { useState, useMemo, useEffect } from "react";
import { Navigate, useParams, useLocation } from "react-router-dom";
import APICaller from "../Requests";

function useQuery() {
  const { search } = useLocation();

  return useMemo(() => new URLSearchParams(search), [search]);
}

function EmailVerification() {
  const [verificationCode, setVerificationCode] = useState("");
  const [message, setMessage] = useState("");
  const [redirect, setRedirect] = useState(false);
  const {userID} = useParams();
  const code = useQuery().get("code");

  useEffect(() => {
    if (code !== "" ) {
      setVerificationCode(code);
    }
  }, [code]);

  const submit = async (e) => {
    e.preventDefault();

    let result;

    try {
      result = await APICaller.VerifyEmail(userID, verificationCode);
    } catch(err) {
      setMessage(err.message);
      return;
    }

    if (result.status !== 200) {
      setMessage(result.data.err);
      return
    } 
    setRedirect(true);
  }

  if (redirect) {
    return <Navigate to="/login"/>;
  }

  return (
    <div className="container mt-4 pt-4">
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
          <form onSubmit={submit}>
            <div className="display-1 mb-4 text-center text-primary"> Verify your email</div>
            <div className="mb-3 text-center text-danger">{message}</div>
            <div className="mb-3 text-center">
              <label className="form-label">Verification Code:</label>
              <input name="code" type="text" className="form-control" id="verification-code" onChange={(e) => setVerificationCode(e.target.value)} value={code}/>
            </div>
            <div className="text-center">
              <button type="submit" className="btn btn-primary text-center">Verify Email</button>
            </div>
          </form>
        </div>
      </div>
      </div>
    );
  }

export default EmailVerification;