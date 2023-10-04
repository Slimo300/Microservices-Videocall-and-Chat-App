import React, {useContext, useState} from 'react';
import {v4 as uuidv4} from "uuid";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import {SetRights, DeleteMember} from '../../requests/Groups';
import { actionTypes, StorageContext } from '../../ChatStorage';
import { UserPicture } from '../Pictures';

export const ModalMembers = ({ group, member, show, toggle }) => {

    const [msg, setMsg] = useState("");

    let message = null;
    if (msg !== "") {
        message = <h5 className='text-danger'>{msg}</h5>
    }

    let nogroup = false;
    if (group.Members === null) {
        nogroup = true
    }
    return (
        <Modal id="buy" tabIndex="-1" size='lg' role="dialog" isOpen={show} toggle={toggle}>
            <div role="document">
                <ModalHeader toggle={toggle} className="bg-dark text-primary text-center">
                    Group Members
                </ModalHeader>
                <ModalBody>
                    <div>
                        {message}
                        <div className='form-group'>
                            <table className="table">
                                <tbody>
                                    {nogroup?null:group.Members.map((item) => {return <Member key={uuidv4()} group={group.ID} member={item} setMsg={setMsg} toggle={toggle} user={member}/>})}
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


const Member = ({ member, group, user, setMsg, toggle }) => {
    const [, dispatch] = useContext(StorageContext);

    const toggleCollapse = () => {
        let elem = document.getElementById("collapse-"+member.ID);
        let isCollapsed = elem.classList.contains("show");
        if (isCollapsed) elem.classList.remove("show")
        else elem.classList.add("show")
    }

    const deleteMember = async() => {

        let response = await DeleteMember(member.groupID, member.ID);
        if (response.status === 200) {
            setMsg("Member deleted");
        } else setMsg(response.data.err);
        setTimeout(function() {
            toggle();
            setMsg("");
        }, 2000);
        dispatch({type: actionTypes.DELETE_MEMBER, payload: {ID: member.ID, groupID: group}})
    }

    return (
        <tr className='d-flex flex-column'>
            <td className="chat-avatar d-flex flex-row justify-content-center">
                <div className='pr-3 members-image-holder'>
                    <UserPicture pictureUrl={member.User.pictureUrl} />
                </div>
                <div className="chat-name pr-3 w-50 d-flex align-items-center">{member.User.username}</div>
                {isSetter(user)?<div className='pr-3 align-right'><button className='btn-primary btn' type="button" onClick={toggleCollapse} disabled={!CanSet(user, member)}>Set rights</button></div>:null}
                {isDeleter(user)?<div className='pr-3 align-right'><button className='btn-primary btn' disabled={!CanDelete(user, member)} onClick={deleteMember}>Delete</button></div>:null}
            </td>
            <Rights member={member} user={user} setMsg={setMsg} />
        </tr>
    );
};

const Rights = ({ member, user, setMsg }) => {

    const [adding, setAdding] = useState(member.adding);
    const toggleAdding = () => {
        setAdding(!adding);
    }
    const [deletingMembers, setDeletingMembers] = useState(member.deletingMembers);
    const toggleDeletingMembers = () => {
        setDeletingMembers(!deletingMembers);
    }
    const [deletingMessages, setDeletingMessages] = useState(member.deletingMembers);
    const toggleDeletingMessages = () => {
        setDeletingMessages(!deletingMessages);
    }
    const [admin, setAdmin] = useState(member.admin);
    const toggleAdmin = () => {
        setAdmin(!admin);
    }

    const setRights = async() => {
        if (adding === member.adding && deletingMembers === member.deleting && admin === member.admin && deletingMessages === member.deletingMessages) {
            return
        }
        let response = await SetRights(
            member.groupID, 
            member.ID, 
            DetermineAction(member.adding, adding), 
            DetermineAction(member.deletingMessages, deletingMessages), 
            DetermineAction(member.deletingMembers, deletingMembers),
            DetermineAction(member.admin, admin)
        );
        console.log(response);
        if (response.status === 200) {
            setMsg("Rights changed");
        } else setMsg(response.data.err);

        setTimeout(function() {
            setMsg("");
        }, 2000);
    }

    return (
        <td className="collapse" id={"collapse-"+member.ID}>
            <div className="card card-body d-flex flex-row">
                <div className='pl-3 d-flex flex-column w-50'>
                    {user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox1" checked={adding} disabled={member.creator} onChange={toggleAdding}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox1">Adding</label>
                    </div>:null}
                    {user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox2" checked={deletingMembers} disabled={member.creator} onChange={toggleDeletingMembers}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox2">Deleting Members</label>
                    </div>:null}
                    {user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox2" checked={deletingMessages} disabled={member.creator} onChange={toggleDeletingMessages}/>
                        <label className="form-check-label" htmlFor="inlineCheckbox2">Deleting Messages</label>
                    </div>:null}
                    {user.admin?<div className='align-middle'>
                        <input className="form-check-input" type="checkbox" id="inlineCheckbox3" checked={admin} disabled={member.creator} onChange={toggleAdmin}/>
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