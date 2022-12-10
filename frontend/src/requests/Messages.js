import axiosObject, {messageService} from "./Setup";


export async function LoadMessages(groupID, offset) {
    return await axiosObject.get(messageService+"/group/"+groupID+"/messages?num=8&offset="+offset);
}

export async function DeleteMessageForYourself(groupID, messageID) {
    return await axiosObject.patch(messageService+"/group/"+groupID+"/messages/"+messageID);
}

export async function DeleteMessageForEveryone(groupID, messageID) {
    return await axiosObject.delete(messageService+"/group/"+groupID+"/messages/"+messageID);
}