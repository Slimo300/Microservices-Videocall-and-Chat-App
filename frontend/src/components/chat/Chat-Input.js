import React, { useRef, useState, useEffect } from "react";
import { GetPresignedRequests } from "../../requests/Messages";

const ChatInput = (props) => {

    const [msg, setMsg] = useState("");
    const [files, setFiles] = useState({});
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
      
    const handleDragOver = (e) => {
        e.preventDefault();
        e.stopPropagation();
    };
      
    const handleDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();  
        const {files} = e.dataTransfer;

        if (files && files.length) {
          let fileInput = document.getElementById("fileUpload")
          fileInput.removeAttribute("hidden");
          fileInput.files = files;
          setFiles(files);
        }
    };

    const sendMessage = async(e) => {
        e.preventDefault();
        if (msg.trim() === "" && files.length === 0) return;

        let fileInfo = [];
        if (files.length !== undefined && files.length > 0) {
            try {
                let res = await GetPresignedRequests(props.group.ID, files.length);
                
                let promises = [];
                for (let i = 0; i < res.data.requests.length; i++) {
                    fileInfo.push({"key": res.data.requests[i].key, "ext": files[i].type})
                    promises.push(fetch(res.data.requests[i].url, {
                        method: 'PUT',
                        body: files[i]
                    }))
                }

                let fileInput = document.getElementById("fileUpload");
                fileInput.setAttribute("hidden", true);
                fileInput.files = null;
                setFiles({});
                await Promise.all(promises);
            } catch(err) {
                alert(err);
            }
        }

        if (props.ws !== undefined) props.ws.send(JSON.stringify({
            "groupID": props.group.ID,
            "userID": props.user.ID,
            "text": msg,
            "nick": props.user.username,
            "files": fileInfo
        }));
        document.getElementById("chat-input").reset();
        document.getElementById("text-area").focus();
    }

    return (
        <form ref={drop} id="chat-input" className="form-group mb-0" onSubmit={sendMessage}>
            <input className="form-control form-control-sm" id="fileUpload" type="file" hidden multiple accept=".jpg, .png, .jpeg"/>
            <div className="d-flex column justify-content-center">
                <textarea autoFocus  id="text-area" className="form-control mr-1" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                <input className="btn btn-primary" type="submit" value="Send"/>
            </div>
        </form>
    )
}

export default ChatInput;