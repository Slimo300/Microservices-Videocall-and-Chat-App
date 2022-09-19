import React from "react";

const GroupMenu = (props) => {

    return (
        <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <button className="dropdown-item" onClick={props.toggleOptions} disabled={!props.member.setting}>Options</button>
            <button className="dropdown-item" onClick={props.toggleMembers}>Members</button>
            <button className="dropdown-item" onClick={props.toggleAdd} disabled={!props.member.adding}>Add User</button>
            <div className="dropdown-divider"></div>
            {props.member.creator?
            <button className="dropdown-item" onClick={props.toggleDel}>Delete Group</button>:
            <button className="dropdown-item" onClick={props.toggleLeave}>Leave Group</button>}
        </div>
    );
};

export default GroupMenu;