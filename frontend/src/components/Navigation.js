import React, { useContext } from 'react';
import {NavLink} from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBell } from '@fortawesome/free-solid-svg-icons'
import {Logout} from '../requests/Users';
import { StorageContext, actionTypes } from '../ChatStorage';
import Invite from './Invite';

const Navigation = (props) => {

    const [state, dispatch] = useContext(StorageContext);

    const logout = async () => {
        let response = await Logout();
        window.localStorage.clear();
        dispatch({type: actionTypes.LOGOUT});
        props.ws.close();
        if (response.status !== 200) alert(response.data.message);
    };

    let menu;

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
                    <button type='button' className="navbar-brand order-1 btn btn-dark text-primary" onClick={props.toggleProfile}>Profile</button>
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
            <NavLink className="navbar-brand" to="/" >ChatApp</NavLink>
            <div className="collapse navbar-collapse" id="navbarCollapse">
                {menu}
            </div>
        </nav>
    )
}

export default Navigation;