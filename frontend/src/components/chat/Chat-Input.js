import React, { useRef, useState, useEffect } from "react";
import { GetPresignedPutRequests } from "../../requests/Messages";

const ChatInput = ({ group, ws, user, member }) => {

    const [msg, setMsg] = useState("");
    const [files, setFiles] = useState(null);

    const form = useRef(null);
    const submitButton = useRef(null);

    const fileInput = useRef(null);
    const textInput = useRef(null);

    useEffect(() => {
        let current = form.current;
        current.addEventListener('dragover', handleDragOver);
        current.addEventListener('drop', handleDrop);
        current.addEventListener('keypress', handleKeypress);
        
        return () => {
            current.removeEventListener('dragover', handleDragOver);
            current.removeEventListener('drop', handleDrop);
            current.removeEventListener('keypress', handleKeypress);
        };
    }, []);
      
    const handleDragOver = (e) => {
        e.preventDefault();
        e.stopPropagation();
    };

    const handleKeypress = (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            e.stopPropagation();
            submitButton.current.click();
        }
    }
      
    const handleDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();  
        const transferFiles = e.dataTransfer.files;

        if (transferFiles && transferFiles.length > 0) {
          fileInput.current.removeAttribute("hidden");
          setFiles(currentFiles => {
                if (!currentFiles) {
                    fileInput.current.files = transferFiles;
                    return transferFiles;
                } else {
                    let newFiles = currentFiles;
                    newFiles.push(...transferFiles);

                    fileInput.current.files = newFiles;
                    return newFiles;
                }
            });
        }
    };

    const sendMessage = async(e) => {
        e.preventDefault();
        if (msg.trim() === "" && files.length === 0) return;

        // fileData stores data about files in message that will be saved to database
        let filesData = [];
        if (files && files.length && files.length > 0) {

            // filesInfo stores data about files that will be needed to generate a presigned request 
            let filesInfo = [];

            for (let i = 0; i < files.length; i++) {
                filesInfo.push({
                    name: `${files[i].name}-${files[i].lastModified}`,
                    size: files[i].size,
                });
            }

            try {
                const response = await GetPresignedPutRequests(group.ID, filesInfo);

                let promises = [];
                for (let i = 0; i < response.data.length; i++) {
                    filesData.push({"key": response.data[i].key, "ext": files[i].type});
                    promises.push(fetch(response.data[i].url, {
                        method: 'PUT',
                        body: files[i],
                    }));
                }
    
                fileInput.current.setAttribute("hidden", true);
                fileInput.current.files = null;
                setFiles(null);
                
                await Promise.all(promises);
            } catch(err) {
                console.log(err);
                // alert(err.response.data.err);
            }
        }

        if (ws !== undefined) ws.send(JSON.stringify({
            Member: {
                ID: member.ID,
                groupID: group.ID,
                userID: user.ID,
                username: user.username,
            },
            text: msg,
            files: filesData
        }));

        form.current.reset();
        
        setMsg("");

        textInput.current.focus();
    }

    return (
        <form ref={form} className="form-group mb-0" onSubmit={sendMessage}>
            <input className="form-control form-control-sm" ref={fileInput} type="file" hidden multiple accept=".jpg, .png, .jpeg" />
            <div className="d-flex column justify-content-center">
                <textarea autoFocus ref={textInput} className="form-control mr-1" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                <input className="btn btn-primary" type="submit" value="Send" ref={submitButton} hidden/>
            </div>
        </form>
    )
}

export default ChatInput;