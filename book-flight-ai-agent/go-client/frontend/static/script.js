// This file originally cloned from https://github.com/yotam-halperin/chatbot-static-UI
const chatbox = document.querySelector(".chatbot .chatbox");
const recordbox = document.querySelector(".record .chatbox");
const chatInput = document.querySelector(".chat-input textarea");
const sendChatBtn = document.querySelector(".chat-input #send-btn");

let userMessage = null; // Variable to store user's message
let userBin = null; // Variable to store user's message
const inputInitHeight = chatInput.scrollHeight;

let fileBlobArr = [];
let fileArr = [];

const createChatLi = (content, className, targetBox = chatbox) => {
    const chatLi = document.createElement("li");
    chatLi.classList.add("chat", `${className}`);

    const iconSpan = document.createElement('span'); // Create icon span here
    if (!(className === "outgoing")) {
        iconSpan.className = "material-symbols-outlined";
        iconSpan.innerText = "smart_toy";
        chatLi.appendChild(iconSpan);
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

    targetBox.appendChild(chatLi);
    targetBox.scrollTo(0, targetBox.scrollHeight);
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
    sendChatBtn.style.display = "none";
    chatInput.style.height = `${inputInitHeight}px`;
    clear();

    // user's message
    createChatLi(contents, "outgoing", chatbox);

    // "Thinking..."
    const incomingChatLi = createChatLi("Thinking...", "incoming", chatbox);
    const incomingRecordLi = createChatLi("Thinking...", "incoming", recordbox); // Add to recordbox

    // timeout
    const TIMEOUT_MS = CONFIG.TIME_OUT_SECOND;
    let isTimeout = false;
    const timeoutId = setTimeout(() => {
        isTimeout = true;
        incomingRecordLi.querySelector("p").textContent = "Request timed out. Please try again.";
    }, TIMEOUT_MS);

    // send request
    generateResponse(incomingChatLi, incomingRecordLi, () => {
        if (!isTimeout) clearTimeout(timeoutId);
    });
}

const generateResponse = (chatElement, recordElement, callback) => {
    const API_URL = "/api/chat";
    const chatMessageElement = chatElement.querySelector("p");
    const recordMessageElement = recordElement.querySelector("p");

    // Initialize stream
    let accumulatedChatResponse = "";
    let accumulatedRecordResponse = "";
    chatMessageElement.textContent = "";
    chatMessageElement.id = "chatContent";
    recordMessageElement.textContent = "";
    recordMessageElement.id = "recordContent";

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
                                if (data.content) {
                                    accumulatedChatResponse += data.content;
                                    chatMessageElement.innerHTML = marked.parse(styleMatch(accumulatedChatResponse)); // Render Markdown
                                    chatbox.scrollTo(0, chatbox.scrollHeight);
                                }
                                if (data.record) {
                                    accumulatedRecordResponse += data.record;
                                    recordMessageElement.innerHTML = marked.parse(styleMatch(accumulatedRecordResponse)); // Render Markdown
                                    recordbox.scrollTo(0, recordbox.scrollHeight);
                                }

                                hljs.highlightAll();
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
        chatMessageElement.classList.add("error");
        chatMessageElement.textContent = "Oops! Something went wrong. Please try again.";
        recordMessageElement.classList.add("error");
        recordMessageElement.textContent = "Oops! Something went wrong. Please try again.";

        // Invoke the callback in case of error
        if (callback) callback();
    });
};

// 渲染航班信息的函数,需要在generateResponse函数中当机票预定成功调用renderFlightInfo渲染，这里假如返回data.flightinfo
const renderFlightInfo = (flightInfo) => {
    const flightInfoContainer = document.getElementById('flight-info');
    // 清空提示信息
    flightInfoContainer.innerHTML = `<h3>航班信息</h3><p>正在加载航班信息...</p>`;
    // 返回为空
    if (flightInfos.length === 0) {
        flightInfoContainer.innerHTML += "<p>没有航班信息可显示</p>";
        return;
    }
    flightInfos.forEach(flightInfo => {
        const flightInfoHTML = `
            <h3>航班信息</h3>
            <p>航班号: ${flightInfo.flightNumber}</p>
            <p>乘客姓名: ${flightInfo.passengerName}</p>
            <p>出发城市: ${flightInfo.departureCity}</p>
            <p>到达城市: ${flightInfo.arrivalCity}</p>
            <p>出发时间: ${flightInfo.departureTime}</p>
            <p>到达时间: ${flightInfo.arrivalTime}</p>
            <hr />
        `;
        flightInfoContainer.innerHTML += flightInfoHTML;
    });
};

chatInput.addEventListener("input", () => {
    if (chatInput.value.trim() !== "") {
        sendChatBtn.style.display = "block";
    } else {
        sendChatBtn.style.display = "none";
    }
    // Adjust the height of the input textarea based on its content
    chatInput.style.height = `${inputInitHeight}px`;
    chatInput.style.height = `${chatInput.scrollHeight}px`;
});
chatInput.addEventListener("keydown", (e) => {
    // If Enter key is pressed without Shift key and the window
    // width is greater than 800px, handle the chat
    if (e.key === "Enter" && !e.shiftKey && window.innerWidth > 800) {
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
        switch (reader.error.code) {
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

// content
function styleMatch(content) {
    // <think>
    // code / equation
    const regex = /(```[\s\S]*?```)|(\$\$[\s\S]*?\$\$)/g;
    return addStringsAroundMatch(content.replace(/\n{2,}/g, '\n').trim(), regex, '\n', '\n');
}

function addStringsAroundMatch(text, regex, prefix, suffix) {
    return text.replace(regex, (match) => {
        return `${prefix}${match}${suffix}`;
    });
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