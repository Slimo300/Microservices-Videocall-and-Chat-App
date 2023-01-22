import React, { useRef, useState, useEffect } from "react";

const ChatInput = (props) => {

    const [msg, setMsg] = useState("");
    const drop = useRef(null);

    useEffect(() => {
        let current = drop.current;
        current.addEventListener('dragover', handleDragOver);
        current.addEventListener('drop', handleDrop);
        
        return () => {
            current.removeEventListener('dragover', handleDragOver);
            current.removeEventListener('drop', handleDrop);
        };
    }, []);
      
    console.log("reload");
    const handleDragOver = (e) => {
        e.preventDefault();
        e.stopPropagation();
    };
      
    const handleDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();  
        const {files} = e.dataTransfer;

        if (files && files.length) {
          console.log(files);
          let fileInput = document.getElementById("fileUpload")
          fileInput.removeAttribute("hidden");
          fileInput.files = files;
        }
    };

    const sendMessage = (e) => {
        e.preventDefault();
        if (msg.trim() === "") return false;
        if (props.ws !== undefined) props.ws.send(JSON.stringify({
            "groupID": props.group.ID,
            "userID": props.user.ID,
            "text": msg,
            "nick": props.user.username,
        }));
        document.getElementById("text-area").value = "";
        document.getElementById("text-area").focus();
    }

    return (
        <form ref={drop} id="chatbox" className="form-group mb-0" onSubmit={sendMessage}>
            <input className="form-control form-control-sm" id="fileUpload" type="file" hidden multiple accept=".jpg, .png, .jpeg"/>
            <div className="d-flex column justify-content-center">
                <textarea autoFocus  id="text-area" className="form-control mr-1" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                <input className="btn btn-primary" type="submit" value="Send" />
            </div>
        </form>
    )
}

export default ChatInput;