// This file originally cloned from https://github.com/yotam-halperin/chatbot-static-UI

const chatbox = document.querySelector(".chatbox");
const chatInput = document.querySelector(".chat-input textarea");
const sendChatBtn = document.querySelector(".chat-input span");

let userMessage = null; // Variable to store user's message
let userBin = null; // Variable to store user's message
const inputInitHeight = chatInput.scrollHeight;

let fileBlobArr = [];
let fileArr = [];

const createChatLi = (message, className) => {
    // Create a chat <li> element with passed message and className
    const chatLi = document.createElement("li");
    chatLi.classList.add("chat", `${className}`);
    let chatContent = className === "outgoing" ? `<p></p>` : `<span class="material-symbols-outlined">smart_toy</span><p></p>`;
    chatLi.innerHTML = chatContent;
    chatLi.querySelector("p").textContent = message;
    return chatLi; // return chat <li> element
}

const handleChat = () => {
    userMessage = chatInput.value.trim(); // Get user entered message and remove extra whitespace
    userBin = null
    if (fileBlobArr.length > 0) {
        userBin = fileBlobArr[0]
    }
    if(!userMessage && !userBin) return;

    // Clear the input textarea and set its height to default
    chatInput.value = "";
    chatInput.style.height = `${inputInitHeight}px`;

    // Append the user's message to the chatbox
    chatbox.appendChild(createChatLi(userMessage, "outgoing"));
    chatbox.scrollTo(0, chatbox.scrollHeight);

    setTimeout(() => {
        // Display "Thinking..." message while waiting for the response
        const incomingChatLi = createChatLi("Thinking...", "incoming");
        chatbox.appendChild(incomingChatLi);
        chatbox.scrollTo(0, chatbox.scrollHeight);
        generateResponse(incomingChatLi);
    }, 600);

    clear()
}

const generateResponse = (chatElement) => {
    const API_URL = "/api/chat";
    const messageElement = chatElement.querySelector("p");

    // init stream
    let accumulatedResponse = "";
    messageElement.textContent = "";
    messageElement.id = "content"

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

            function readStream() {
                return reader.read().then(({ done, value }) => {
                    if (done) {
                        return;
                    }

                    const chunk = decoder.decode(value);
                    const events = chunk.split('\n\n');

                    events.forEach(event => {
                        if (event.startsWith('event:message')) {
                            // extract data
                            const dataLine = event.split('\n')[1];
                            if (dataLine && dataLine.startsWith('data:')) {
                                try {
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

                    return readStream();
                });
            }

            return readStream();
        })
        .catch(error => {
            console.error('Error:', error);
            messageElement.classList.add("error");
            messageElement.textContent = "Oops! Something went wrong. Please try again.";
        });
}

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
        img.src = "../static/file.png";

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