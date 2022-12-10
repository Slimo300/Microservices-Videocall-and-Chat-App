import React, {useContext, useState} from 'react';
import {v4 as uuidv4} from "uuid";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import {SetRights, DeleteMember} from '../../requests/Groups';
import { actionTypes, StorageContext } from '../../ChatStorage';
import { UserPicture } from '../Pictures';

export const ModalMembers = (props) => {

    const [msg, setMsg] = useState("");

    let message = null;
    if (msg !== "") {
        message = <h5 className='text-danger'>{msg}</h5>
    }

    let nogroup = false;
    if (props.group.Members === null) {
        nogroup = true
    }
    return (
        <Modal id="buy" tabIndex="-1" size='lg' role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Group Members
                </ModalHeader>
                <ModalBody>
                    <div>
                        {message}
                        <div className='form-group'>
                            <table className="table">
                                <tbody>
                                    {nogroup?null:props.group.Members.map((item) => {return <Member key={uuidv4()} group={props.group.ID} member={item} setMsg={setMsg} toggle={props.toggle} user={props.member}/>})}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 

function isDeleter(member) {
    if (member.creator || member.admin || member.deletingMembers) return true;
    return false;
}

function isSetter(member) {
    if (member.creator || member.admin) return true;
    return false;
}

function CanDelete(issuer, target) {
    if (target.creator) return false;
    if (target.admin && !issuer.creator) return false;
    if (target.deletingMembers && !issuer.creator && !issuer.admin) return false;
    if (!issuer.deletingMembers) return false;
    return true;
}

function CanSet(issuer, target) {
    if (target.creator) return false;
    if (target.admin && !issuer.creator) return false;
    if (!issuer.admin) return false;
    return true;
}


const Member = (props) => {
    const [, dispatch] = useContext(StorageContext);

    const toggleCollapse = () => {
        let elem = document.getElementById("collapse-"+props.member.ID);
        let isCollapsed = elem.classList.contains("show");
        if (isCollapsed) elem.classList.remove("show")
        else elem.classList.add("show")
    }

    const deleteMember = async() => {

        let response = await DeleteMember(props.member.groupID, props.member.ID);
        if (response.status === 200) {
            props.setMsg("Member deleted");
        } else props.setMsg(response.data.err);
        setTimeout(function() {
            props.toggle();
            props.setMsg("");
        }, 2000);
        dispatch({type: actionTypes.DELETE_MEMBER, payload: {ID: props.member.ID, groupID: props.group}})
    }

    return (
        <tr className='d-flex flex-column'>
            <td className="chat-avatar d-flex flex-row justify-content-center">
                <div className='pr-3 members-image-holder'>
                    <UserPicture pictureUrl={props.member.User.pictureUrl} />
                </div>
                <div className="chat-name pr-3 w-50 d-flex align-items-center">{props.member.User.username}</div>
                {isSetter(props.user)?<div className='pr-3 align-right'><button className='btn-primary btn' type="button" onClick={toggleCollapse} disabled={!CanSet(props.user, props.member)}>Set rights</button></div>:null}
                {isDeleter(props.user)?<div className='pr-3 align-right'><button className='btn-primary btn' disabled={!CanDelete(props.user, props.member)} onClick={deleteMember}>Delete</button></div>:null}
            </td>
            <Rights member={props.member} user={props.user} setMsg={props.setMsg} />
        </tr>
    );
};

const Rights = (props) => {

    const [adding, setAdding] = useState(props.member.adding);
    const toggleAdding = () => {
        setAdding(!adding);
    }
    const [deletingMembers, setDeletingMembers] = useState(props.member.deletingMembers);
    const toggleDeletingMembers = () => {
        setDeletingMembers(!deletingMembers);
    }
    const [deletingMessages, setDeletingMessages] = useState(props.member.deletingMembers);
    const toggleDeletingMessages = () => {
        setDeletingMessages(!deletingMessages);
    }
    const [admin, setAdmin] = useState(props.member.admin);
    const toggleAdmin = () => {
        setAdmin(!admin);
    }

    const setRights = async() => {
        if (adding === props.member.adding && deletingMembers === props.member.deleting && admin === props.member.admin && deletingMessages === props.member.deletingMessages) {
            return
        }
        let response = await SetRights(
            props.member.groupID, 
            props.member.ID, 
            DetermineAction(props.member.adding, adding), 
            DetermineAction(props.member.deletingMessages, deletingMessages), 
            DetermineAction(props.member.deletingMembers, deletingMembers),
            DetermineAction(props.member.admin, admin)
        );
        console.log(response);
        if (response.status === 200) {
            props.setMsg("Rights changed");
        } else props.setMsg(response.data.err);

        setTimeout(function() {
            props.setMsg("");
        }, 2000);
    }

    return (
        <td className="collapse" id={"collapse-"+props.member.ID}>
            <div className="card card-body d-flex flex-row">
                <div className='pl-3 d-flex flex-column w-50'>
                    {props.user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox1" checked={adding} disabled={props.member.creator} onChange={toggleAdding}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox1">Adding</label>
                    </div>:null}
                    {props.user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox2" checked={deletingMembers} disabled={props.member.creator} onChange={toggleDeletingMembers}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox2">Deleting Members</label>
                    </div>:null}
                    {props.user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox2" checked={deletingMessages} disabled={props.member.creator} onChange={toggleDeletingMessages}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox2">Deleting Messages</label>
                    </div>:null}
                    {props.user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox3" checked={admin} disabled={props.member.creator} onChange={toggleAdmin}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox3">Admin</label>
                    </div>:null}
                </div>
                <div className='d-flex flex-row justify-content-end align-items-center w-50'>
                    <button className='btn btn-secondary' onClick={setRights}>Change Rights</button>
                </div>
            </div>
        </td>
    )
}

const DetermineAction = (originalState, submitState) => {
    if (originalState === submitState) return 0;
    if (originalState === false && submitState === true) return 1;
    if (originalState === true && submitState === false) return -1;
    throw new Error("State are not booleans");
}