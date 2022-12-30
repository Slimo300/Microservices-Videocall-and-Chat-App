import axiosObject, {groupsService} from "./Setup";

export async function GetInvites(offset) {
    return await axiosObject.get(groupsService+"/invites?num=8&offset="+offset);
}

export async function GetGroups() {
    return await axiosObject.get(groupsService+"/group");
}

export async function CreateGroup(name) {
    return await axiosObject.post(groupsService+"/group", {
        "name": name,
    })
}

export async function DeleteGroup(groupID) {
    return await axiosObject.delete(groupsService+"/group/"+groupID);
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

export async function DeleteMember(groupID, memberID) {
    return await axiosObject.delete(groupsService+"/group/"+groupID+"/member/"+memberID);
}

export async function SetRights(groupID, memberID, adding, deletingMessages, deletingMembers, setting) {
    return await axiosObject.patch(groupsService+"/group/"+groupID+"/member/"+memberID, {
        "adding": adding,
        "deletingMessages": deletingMessages,
        "deletingMembers": deletingMembers,
        "setting": setting
    });
}

export async function UpdateGroupProfilePicture(imageForm, groupID) {
    return await axiosObject.post(groupsService+"/group/"+groupID+"/image", imageForm, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    });
}

export async function DeleteGroupProfilePicture(groupID) {
    return await axiosObject.delete(groupsService+"group/"+groupID+"/image");
}
