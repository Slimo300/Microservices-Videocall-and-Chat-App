import React from "react";

const GroupMenu = ({ toggleOptions, toggleMembers, toggleAdd, toggleDel, toggleLeave, member}) => {

    return (
        <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <button className="dropdown-item" onClick={toggleOptions} disabled={!member.admin}>Options</button>
            <button className="dropdown-item" onClick={toggleMembers}>Members</button>
            <button className="dropdown-item" onClick={toggleAdd} disabled={!member.adding}>Add User</button>
            <div className="dropdown-divider"></div>
            {member.creator?
            <button className="dropdown-item" onClick={toggleDel}>Delete Group</button>:
            <button className="dropdown-item" onClick={toggleLeave}>Leave Group</button>}
        </div>
    );
};

export default GroupMenu;