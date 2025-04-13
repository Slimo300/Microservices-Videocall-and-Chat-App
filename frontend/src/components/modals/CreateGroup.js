import React, { useContext, useState } from 'react';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { actionTypes, StorageContext } from '../../ChatStorage';
import { CreateGroup } from '../../requests/Groups';
import axiosObject from '../../requests/Setup';

export const ModalCreateGroup = ({ toggle, show }) => {
    const [, dispatch] = useContext(StorageContext);

    const [groupName, setGroupName] = useState("");
    const [msg, setMsg] = useState("");

    const submit = async(e) => {
        e.preventDefault();
        try {
            const createGroupResponse = await CreateGroup(groupName);
            const response = await axiosObject.get(createGroupResponse.headers["content-location"]);
            dispatch({type: actionTypes.ADD_GROUP, payload: response.data});
            setMsg("Group created");

            setTimeout(function () {
                toggle();
                setMsg("");
            }, 1000);
        }
        catch(err) {
            console.log(err)
            if (err.response.data.err !== undefined) setMsg(err.response.data.err);
            else setMsg(err.message);
        }
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={show} toggle={toggle}>
            <div role="document">
                <ModalHeader toggle={toggle} className="bg-dark text-primary text-center">
                    Create Group
                </ModalHeader>
                <ModalBody>
                    <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form onSubmit={submit}>
                            <div className="form-group">
                                <label htmlFor="email">Group name:</label>
                                <input name="name" type="text" className="form-control" onChange={(e)=>{setGroupName(e.target.value)}}/>
                            </div>
                            <div className="form-row text-center">
                                <div className="col-12 mt-2">
                                    <button type="submit" className="btn btn-dark btn-large text-primary">Create group</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 