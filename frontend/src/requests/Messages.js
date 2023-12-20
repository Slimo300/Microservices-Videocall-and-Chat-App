import axiosObject, {messageService} from "./Setup";


export async function LoadMessages(groupID, offset) {
    return await axiosObject.get(messageService+"/"+groupID+"?num=8&offset="+offset);
}

export async function DeleteMessageForYourself(messageID) {
    return await axiosObject.patch(messageService+"/"+messageID+"/hide");
}

export async function DeleteMessageForEveryone(messageID) {
    return await axiosObject.delete(messageService+"/"+messageID);
}

export async function GetPresignedRequests(groupID, files) {
    return await axiosObject.post(messageService+"/"+groupID+"/presign", {
        "files": files,
    });
}