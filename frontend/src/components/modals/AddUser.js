import React, {useContext, useEffect, useState} from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { StorageContext } from "../../ChatStorage";
import { SendGroupInvite } from "../../requests/Groups";
import { actionTypes } from "../../ChatStorage";
import { SearchUsers } from "../../requests/Users";
import { UserPicture } from "../Pictures";

export const ModalAddUser = ({ group, toggle, show }) => {

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
                // setMsg(err.response.data.err);
            }
        };
        fetchData();
    }, [username]);
    
    return (
        <Modal tabIndex="-1" role="dialog" isOpen={show} toggle={toggle}>
            <div role="document">
                <ModalHeader toggle={toggle} className="bg-dark text-primary text-center">
                    Add User
                </ModalHeader>
                <ModalBody>
                    <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form onSubmit={(e) => {e.preventDefault()}}>
                            <div className="form-group">
                                <label htmlFor="email">Username:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" autoComplete="off" onChange={(e)=>{setUsername(e.target.value)}}/>
                            </div>
                            <div className="dropdown w-100">
                            <div data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"/>
                            <div id="dropdownUsers" className="dropdown-menu w-100" aria-labelledby="dropdownMenuButton">
                                {users===null||users.length===0?null:users.map((item) => {
                                    return <div key={item.ID}>
                                            <User user={item} setMsg={setMsg} groupID={group.ID} toggle={toggle} isMember={isMember(group, item.ID)}/>
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


const User = ({ user, groupID, setMsg, toggle, isMember }) => {
    const [, dispatch] = useContext(StorageContext);

    const AddUser = async(e) => {
        e.preventDefault();
        try {
            let response = await SendGroupInvite(user.ID, groupID);
            dispatch({type: actionTypes.ADD_INVITE, payload: response.data});

            setMsg("Invite sent successfully");

            setTimeout(function () {    
                toggle();
                setMsg("");
            }, 1500);

        } catch(err) {
            setMsg(err.response);
        }
    }
    console.log(user);

    return (
        <div className="d-flex column justify-content-between align-items-center px-3">
            <div className="d-flex column align-items-center">
                <div className="chat-avatar image-holder-invite"><UserPicture userID={user.ID} hasPicture={user.hasPicture}/></div>
                <div className="user-name pl-3">{user.username}</div>
            </div>
            <button className="btn btn-primary pl-3" disabled={isMember} onClick={AddUser}>Add User</button>
        </div>
    )
}


function isMember(group, userID) {
    for (let i = 0; i < group.Members.length; i++) {
        if (group.Members[i].userID === userID) return true;
    }
    return false;
}