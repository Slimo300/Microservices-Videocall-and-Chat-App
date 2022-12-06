import React, {useContext, useState} from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { StorageContext } from "../../ChatStorage";
import {SendGroupInvite} from "../../Requests";
import { actionTypes } from "../../ChatStorage";

export const ModalAddUser = (props) => {

    const [username, setUsername] = useState("");
    const [msg, setMsg] = useState("");
    const [, dispatch] = useContext(StorageContext);

    const submitAddToGroup = async(e) => {
        e.preventDefault();

        try {
            let response = await SendGroupInvite(username, props.group.ID);
            console.log(response.data);
            dispatch({type: actionTypes.ADD_INVITE, payload: response.data});
            setMsg("Invite sent successfully");
            setTimeout(function () {    
                props.toggle();
                setMsg("");
            }, 3000);
        } catch(err) {
            setMsg(err.response.data.err);
        }
    }
    
    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Add User
                </ModalHeader>
                <ModalBody>
                    <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form onSubmit={submitAddToGroup}>
                            <div className="form-group">
                                <label htmlFor="email">Username:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" onChange={(e)=>{setUsername(e.target.value)}}/>
                            </div>
                            <div className="form-row text-center">
                                <div className="col-12 mt-2">
                                    <button type="submit" className="btn btn-dark btn-large text-primary">Add User</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 
