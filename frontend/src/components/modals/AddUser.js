import React, {useContext, useEffect, useState} from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { StorageContext } from "../../ChatStorage";
import {SendGroupInvite} from "../../requests/Groups";
import { actionTypes } from "../../ChatStorage";
import { SearchUsers } from "../../requests/Users";
import { UserPicture } from "../Pictures";

export const ModalAddUser = (props) => {

    const [username, setUsername] = useState("");
    const [msg, setMsg] = useState("");

    const [users, setUsers] = useState([]);

    useEffect( () => {
        async function fetchData() {
            let dropdown = document.getElementById("dropdownUsers");
            if (dropdown === null) return;
            if (username.length < 4) {
                dropdown.classList.remove("show");
                return;
            }
            try {
                let response = await SearchUsers(username, 5);
                setUsers(response.data);
                dropdown.classList.add("show");
            }
            catch(err) {
                setMsg(err);
            }
        };
        fetchData();
    }, [username]);
    
    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Add User
                </ModalHeader>
                <ModalBody>
                    <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form>
                            <div className="form-group">
                                <label htmlFor="email">Username:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" autoComplete="off" onChange={(e)=>{setUsername(e.target.value)}}/>
                            </div>
                            <div className="dropdown w-100">
                            <div data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"/>
                            <div id="dropdownUsers" className="dropdown-menu w-100" aria-labelledby="dropdownMenuButton">
                                {users===null||users.length===0?null:users.map((item) => {
                                    return <div key={item.ID}>
                                            <User user={item} setMsg={setMsg} groupID={props.group.ID} toggle={props.toggle} isMember={isMember(props.group, item.ID)}/>
                                            <hr />
                                        </div>
                                })}
                            </div>      
                        </div>
                        </form>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 


const User = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const AddUser = async(e) => {
        e.preventDefault();
        try {
            let response = await SendGroupInvite(props.user.ID, props.groupID);
            dispatch({type: actionTypes.ADD_INVITE, payload: response.data});

            props.setMsg("Invite sent successfully");

            setTimeout(function () {    
                props.toggle();
                props.setMsg("");
            }, 3000);

        } catch(err) {
            props.setMsg(err.response);
        }
    }

    return (
        <div className="d-flex column justify-content-between align-items-center px-3">
            <div className="d-flex column align-items-center">
                <div className="chat-avatar image-holder-invite"><UserPicture pictureUrl={props.user.picture}/></div>
                <div className="user-name pl-3">{props.user.username}</div>
            </div>
            <button className="btn btn-primary pl-3" disabled={props.isMember} onClick={AddUser}>Add User</button>
        </div>
    )
}


function isMember(group, userID) {
    for (let i = 0; i < group.Members.length; i++) {
        if (group.Members[i].userID === userID) return true;
    }
    return false;
}