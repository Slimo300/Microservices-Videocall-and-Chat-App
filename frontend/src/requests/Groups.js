import axiosObject, {groupsService} from "./Setup";

export async function GetInvites(offset) {
    return await axiosObject.get(groupsService+"/invites?num=8&offset="+offset);
}

export async function GetGroups() {
    return await axiosObject.get(groupsService+"/");
}

export async function CreateGroup(name) {
    return await axiosObject.post(groupsService+"/", {
        "name": name,
    })
}

export async function DeleteGroup(groupID) {
    return await axiosObject.delete(groupsService+"/"+groupID);
}

export async function SendGroupInvite(targetID, groupID) {
    return await axiosObject.post(groupsService+"/invites", {
        "target": targetID,
        "group": groupID
    });
}

export async function RespondGroupInvite(inviteID, answer) {
    return await axiosObject.put(groupsService+"/invites/"+inviteID, {
        "answer": answer
    })
}

export async function DeleteMember(memberID) {
    return await axiosObject.delete(groupsService+"/members/"+memberID);
}

export async function SetRights(memberID, adding, deletingMessages, muting, deletingMembers, admin ) {
    return await axiosObject.patch(groupsService+"/members/"+memberID, {
        "adding": adding,
        "deletingMessages": deletingMessages,
        "deletingMembers": deletingMembers,
        "muting": muting,
        "admin": admin,
    });
}

export async function UpdateGroupProfilePicture(imageForm, groupID) {
    return await axiosObject.post(groupsService+"/"+groupID+"/image", imageForm, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    });
}

export async function DeleteGroupProfilePicture(groupID) {
    return await axiosObject.delete(groupsService+"/"+groupID+"/image");
}
