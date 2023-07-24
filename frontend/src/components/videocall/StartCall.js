import mockVideo from "../../statics/videos/mock.webm";

const StartCall = async (mocking, initVideo) => {

    let stream;

    if (mocking) {
        const video = document.createElement("video");
        video.src = mockVideo;
        video.volume = 0.1;
        stream = video.captureStream();
        await video.play();
    } else {
        stream = await navigator.mediaDevices.getUserMedia({video: initVideo, audio: true});
    }

    return stream;
};

export default StartCall;