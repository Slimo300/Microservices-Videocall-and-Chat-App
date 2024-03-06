import React from "react";

const GroupMenu = ({ toggleOptions, toggleMembers, toggleAdd, toggleDel, toggleLeave, member}) => {

    return (
        <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <button className="dropdown-item" onClick={toggleOptions} disabled={!member.admin && !member.creator}>Options</button>
            <button className="dropdown-item" onClick={toggleMembers}>Members</button>
            <button className="dropdown-item" onClick={toggleAdd} disabled={!member.adding && !member.admin && !member.creator}>Add User</button>
            <div className="dropdown-divider"></div>
            {member.creator?
            <button className="dropdown-item" onClick={toggleDel}>Delete Group</button>:
            <button className="dropdown-item" onClick={toggleLeave}>Leave Group</button>}
        </div>
    );
};

export default GroupMenu;