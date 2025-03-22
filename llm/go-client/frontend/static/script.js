// This file originally cloned from https://github.com/yotam-halperin/chatbot-static-UI

const chatbox = document.querySelector(".chatbox");
const chatInput = document.querySelector(".chat-input textarea");
const sendChatBtn = document.querySelector(".chat-input span");

let userMessage = null; // Variable to store user's message
let userBin = null; // Variable to store user's message
const inputInitHeight = chatInput.scrollHeight;

let fileBlobArr = [];
let fileArr = [];

const createChatLi = (content, className) => {
    const chatLi = document.createElement("li");
    chatLi.classList.add("chat", `${className}`);

    if (!(className === "outgoing")) {
        let toy = document.createElement('span');
        toy.className = "material-symbols-outlined";
        toy.innerText = "smart_toy"
        chatLi.appendChild(toy);
    }

    const contents = Array.isArray(content) ? content : [content];

    contents.forEach(item => {
        if (!item) return;
        if (item.startsWith('data:image')) {
            const img = document.createElement('img');
            img.src = item;
            chatLi.appendChild(img);
        } else {
            const p = document.createElement('p');
            p.textContent = item;
            chatLi.appendChild(p);
        }
    });

    return chatLi;
};

const handleChat = () => {
    userMessage = chatInput.value.trim();
    userBin = fileBlobArr.length > 0 ? fileBlobArr[0] : null;

    const contents = [];
    if (userMessage) contents.push(userMessage);
    if (userBin) contents.push(userBin);

    if (contents.length === 0) return;

    chatInput.value = "";
    chatInput.style.height = `${inputInitHeight}px`;
    clear();

    // user's message
    chatbox.appendChild(createChatLi(contents, "outgoing"));
    chatbox.scrollTo(0, chatbox.scrollHeight);

    // "Thinking..."
    const incomingChatLi = createChatLi("Thinking...", "incoming");
    chatbox.appendChild(incomingChatLi);
    chatbox.scrollTo(0, chatbox.scrollHeight);

    // timeout
    const TIMEOUT_MS = CONFIG.TIME_OUT_SECOND;
    let isTimeout = false;
    const timeoutId = setTimeout(() => {
        isTimeout = true;
        incomingChatLi.querySelector(".content").textContent = "Request timed out. Please try again.";
        chatbox.scrollTo(0, chatbox.scrollHeight);
    }, TIMEOUT_MS);

    // send request
    generateResponse(incomingChatLi, () => {
        if (!isTimeout) clearTimeout(timeoutId);
    });
}

const generateResponse = (chatElement, callback) => {
    const API_URL = "/api/chat";
    const messageElement = chatElement.querySelector("p");

    // Initialize stream
    let accumulatedResponse = "";
    messageElement.textContent = "";
    messageElement.id = "content";

    fetch(API_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ message: userMessage, bin: userBin }),
    })
        .then(response => {
            const reader = response.body.getReader();
            const decoder = new TextDecoder();

            // Function to read the stream recursively
            function readStream() {
                return reader.read().then(({ done, value }) => {
                    if (done) {
                        // Stream is complete, invoke the callback
                        if (callback) callback();
                        return;
                    }

                    // Decode the chunk and process events
                    const chunk = decoder.decode(value);
                    const events = chunk.split('\n\n');

                    events.forEach(event => {
                        if (event.startsWith('event:message')) {
                            // Extract data from the event
                            const dataLine = event.split('\n')[1];
                            if (dataLine && dataLine.startsWith('data:')) {
                                try {
                                    // Parse the JSON data and update the UI
                                    const data = JSON.parse(dataLine.replace('data:', ''));
                                    accumulatedResponse += data.content;
                                    messageElement.textContent = accumulatedResponse;
                                    chatbox.scrollTo(0, chatbox.scrollHeight);
                                } catch (error) {
                                    console.error('Failed to parse event data:', error);
                                }
                            }
                        }
                    });

                    // Continue reading the stream
                    return readStream();
                });
            }

            // Start reading the stream
            return readStream();
        })
        .catch(error => {
            console.error('Error:', error);
            messageElement.classList.add("error");
            messageElement.textContent = "Oops! Something went wrong. Please try again.";

            // Invoke the callback in case of error
            if (callback) callback();
        });
};

chatInput.addEventListener("input", () => {
    // Adjust the height of the input textarea based on its content
    chatInput.style.height = `${inputInitHeight}px`;
    chatInput.style.height = `${chatInput.scrollHeight}px`;
});
chatInput.addEventListener("keydown", (e) => {
    // If Enter key is pressed without Shift key and the window
    // width is greater than 800px, handle the chat
    if(e.key === "Enter" && !e.shiftKey && window.innerWidth > 800) {
        e.preventDefault();
        handleChat();
    }
});
sendChatBtn.addEventListener("click", handleChat);

addBtn = document.getElementById("add-btn");
fileInput = document.getElementById("input");

// file process
function filesToBlob(file) {
    let reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = e => {
        fileBlobArr.push(e.target.result);
        let fileDiv = document.createElement('div');
        // delete btn
        let removeDiv = document.createElement('div');
        removeDiv.id = 'file' + '-' + fileBlobArr.length;
        removeDiv.innerHTML = '×';
        // file name
        let fileName = document.createElement('p');
        fileName.innerHTML = file.name;
        fileName.title = file.name;

        let img = document.createElement('img');
        img.src = e.target.result;


        fileDiv.appendChild(img);
        fileDiv.appendChild(removeDiv);
        fileDiv.appendChild(fileName);

        document.getElementById("drop").appendChild(fileDiv);
    };

    reader.onerror = () => {
        switch(reader.error.code) {
            case '1':
                alert('File not found');
                break;
            case '2':
                alert('Security error');
                break;
            case '3':
                alert('Loading interrupted');
                break;
            case '4':
                alert('File is not readable');
                break;
            case '5':
                alert('Encode error');
                break;
            default:
                alert('File read fail');
        }
    };
}

function handleFileSelect(event) {
    const files = event.target.files;
    if (files.length > 0) {
        const file = files[0];

        if (!file.type.startsWith('image/')) {
            alert("Only support image files");
            return;
        }

        fileArr.push(file);
        filesToBlob(file);

        document.querySelector('.drop-box').style.setProperty('--div-count', "1");
        document.getElementById("drop").style.display = "flex";
        addBtn.style.display = "none";
    }
}

fileInput.addEventListener('change', handleFileSelect);

addBtn.addEventListener('click', () => {
    fileInput.click();
    document.getElementById("drop").style.display = "flex";
});

function clear() {
    document.getElementById("drop").innerHTML = '';
    document.getElementById("drop").style.display = "none";
    addBtn.style.display = "flex";
    fileInput.value = "";
    fileBlobArr = [];
    fileArr = [];
}

document.getElementById("drop").addEventListener('click', clear);