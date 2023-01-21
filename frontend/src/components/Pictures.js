import group from "../images/group.jpg";
import user from "../images/user.png";

export const UserPicture = (props) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="user" id={props.imageID}
            src={(props.pictureUrl === "")?user:"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.pictureUrl}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=user;
            }}
        />
    );
} 
export const GroupPicture = (props) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="group" id={props.imageID}
            src={(props.pictureUrl === "")?group:"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.pictureUrl}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=group;
            }}
        />
    );
} 