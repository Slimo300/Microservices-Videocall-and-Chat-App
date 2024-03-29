import React, { useContext, useEffect } from 'react';
import {NavLink} from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBell } from '@fortawesome/free-solid-svg-icons';

import {Logout} from '../requests/Users';
import { StorageContext, actionTypes } from '../ChatStorage';
import Invite from './Invite';
import logo from "../statics/images/logo.png";

const Navigation = ({ws, toggleProfile, setWs}) => {

    const [state, dispatch] = useContext(StorageContext);

    const logout = async () => {
        try {
            await Logout();
        } catch (err) {
        }
        window.localStorage.clear();
        dispatch({type: actionTypes.LOGOUT});
        try {
            ws.close();
        } catch (err) {

        }
        // setWs changes state and triggers nav rerender
        setWs({});
    };

    useEffect(() => {      
        if (dispatch && setWs) {
            window.addEventListener("logout", () => {
                window.localStorage.clear();
                dispatch({ type: actionTypes.LOGOUT });
                setWs(ws => {
                    try {
                        ws.close()
                        return null;
                    } catch (err) {
                        console.log(err);
                    }
                })
            });
        }  
    }, [dispatch, setWs]);

    let menu;

    if (window.location.pathname.match("/call/*")){
        return null;
    }

    if (window.localStorage.getItem("token") === null) {
        menu = (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto"></ul>
                <NavLink className="nav-item nav-link" to="/login">Login</NavLink>
                <NavLink className="nav-item nav-link" to="/register">Register</NavLink>
            </div>
        );
    } else {
        menu = (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto">
                    <button type='button' className="navbar-brand order-1 btn btn-dark text-primary" onClick={toggleProfile}>Profile</button>
                </ul>

                <div className="btn-group">
                    <button type="button" className="btn btn-secondary dropdown-toggle bg-dark" data-toggle="dropdown" aria-expanded="false">
                        <FontAwesomeIcon className='text-primary pr-2' icon={faBell} />
                        <span className="badge badge-secondary">{state.invites!==undefined?state.invites.length:null}</span>
                    </button>
                    <div className="dropdown-menu dropdown-menu-right">
                        {state.invites!==undefined?state.invites.map((item)=> <Invite key={item.ID} invite={item} userID={state.user.ID} />):null}
                    </div>
                </div>

                <NavLink className="nav-item nav-link" to="/login?logout=true" onClick={logout}>Logout</NavLink>
            </div>
        );
    }

    return (
        <nav className="navbar navbar-expand-md navbar-dark bg-dark mb-4">
            <NavLink className="navbar-brand" to="/" >
                <img src={logo} alt="Logo" width="200" height="55" className="d-inline-block align-text-top" />
            </NavLink>
            <div className="collapse navbar-collapse" id="navbarCollapse">
                {menu}
            </div>
        </nav>
    )
}

export default Navigation;