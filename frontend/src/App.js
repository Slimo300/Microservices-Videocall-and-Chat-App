import React, {useState} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";

import Navigation from "./components/Navigation";
import ChatStorage from "./ChatStorage";

import Main from "./pages/Main";
import RegisterForm from "./pages/Register";
import EmailVerification from "./pages/VerifyEmail";
import Page404 from "./pages/404";
import SignInForm from './pages/Login';

function App() {
  
  const [name, setName] = useState('');
  const [ws, setWs] = useState({}); // websocket connection

  const [profileShow, setProfileShow] = useState(false);
  const toggleProfileShow = () => {
    setProfileShow(!profileShow);
  }

  return (
      <div >
        <ChatStorage>
        <Router>
          <Navigation name={name} setName={setName} toggleProfile={toggleProfileShow} ws={ws} />
          <main>
            <Routes>
              <Route path="/" element={<Main name={name} profileShow={profileShow} toggleProfile={toggleProfileShow} ws={ws} setWs={setWs}/>}/>
              <Route path="/login" element={<SignInForm setName={setName} name={name}/>}/>
              <Route path="/register" element={<RegisterForm/>}/>
              <Route path="/verify-email/:userID" element={<EmailVerification />}/>
              <Route path="*" element={<Page404 />}/>
            </Routes>
          </main>
        </Router>
        </ChatStorage>
      </div>
  )
}

export default App;