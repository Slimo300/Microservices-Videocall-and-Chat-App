import React, { useState, useEffect } from "react";
import {Navigate, useParams} from "react-router-dom";
import {VerifyAccount} from "../requests/Users";

function EmailVerification() {
  const [message, setMessage] = useState("");
  const [redirect, setRedirect] = useState(false);
  const {code} = useParams();

  useEffect(() => {

    const verify = async () => {
      try {
        await VerifyAccount(code);
      } catch(err) {
        setMessage(err.response.data.err);
        setRedirect(true);
        return;
      }

      setMessage("Account activated");
      setRedirect(true);
    }

    if (code !== "" ) {
      verify()
    }
  }, [code]);

  if (redirect) {
    return <Navigate to={"/login?message="+message.replaceAll(" ", "+")}/>;
  }

  return null;
}

export default EmailVerification;