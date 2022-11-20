import React, {useContext, useState} from 'react';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { actionTypes, StorageContext } from '../../ChatStorage';
import {CreateGroup} from '../../Requests';

export const ModalCreateGroup = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const [grName, setGrName] = useState("");
    const [msg, setMsg] = useState("");

    const submit = async(e) => {
        e.preventDefault();
        let group = await CreateGroup(grName);
        if (group.data.name === "Error") {
            setMsg(group.data);
        } else {
            dispatch({type: actionTypes.NEW_GROUP, payload: group.data});
            setMsg("Group created");
        }
        setTimeout(function () {    
            props.toggle();
            setMsg("");
        }, 1000);
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Create Group
                </ModalHeader>
                <ModalBody>
                    <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form onSubmit={submit}>
                            <div className="form-group">
                                <label htmlFor="email">Group name:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" onChange={(e)=>{setGrName(e.target.value)}}/>
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