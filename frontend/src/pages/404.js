import React from 'react';
import { Navigate } from "react-router-dom";

function Page404() {
  if (window.localStorage.getItem("token") !== null) return <Navigate to="/"/>

  return (
    <div className="container mt-4 pt-4">
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
        <div className="display-1 mb-4 text-center text-primary"><div className='text-danger'>404</div>Page not found</div>
        </div>
      </div>
      </div>
    );
  }

export default Page404;