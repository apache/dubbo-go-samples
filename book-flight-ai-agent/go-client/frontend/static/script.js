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

// 1. 页面初始化时请求后端配置（新增）
async function loadConfig() {
    try {
        const res = await fetch("/api/config"); // 后端新增接口返回配置
        window.CONFIG = await res.json();
        // 统一超时字段名：使用与后端一致的 TIMEOUT_SECONDS（原 TIME_OUT_SECOND 废弃）
        window.CONFIG.TIMEOUT_MS = window.CONFIG.TIMEOUT_SECONDS * 1000; // 转为毫秒（方便定时器使用）
    } catch (err) {
        console.error("加载配置失败，使用默认超时（2分钟）", err);
        window.CONFIG = window.CONFIG || {};
        window.CONFIG.TIMEOUT_SECONDS = 120;
        window.CONFIG.TIMEOUT_MS = 120 * 1000;
    }
}

// 2. 页面加载完成后执行配置加载
window.onload = async () => {
    // 确保window.CONFIG已经存在，如果不存在则初始化
    if (!window.CONFIG) {
        await loadConfig();
    } else {
        // 确保window.CONFIG有TIMEOUT_MS属性
        if (!window.CONFIG.TIMEOUT_MS && window.CONFIG.TIME_OUT_SECOND) {
            window.CONFIG.TIMEOUT_MS = window.CONFIG.TIME_OUT_SECOND;
        }
    }
    // 其他初始化逻辑（如绑定按钮事件）
};

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


    // 超时逻辑优化
    let timeoutId;
    const startTimeout = () => {
        // 清除已有定时器（避免重复）
        if (timeoutId) clearTimeout(timeoutId);
        // 启动新定时器（使用同步后的 window.CONFIG.TIMEOUT_MS）
        timeoutId = setTimeout(() => {
            const timeoutMsg = `
        <div>
          <p>请求超时（当前超时时间：${window.CONFIG.TIMEOUT_SECONDS || window.CONFIG.TIME_OUT_SECOND/1000}秒）</p>
          <button class="retry-btn" style="margin-top:8px;padding:4px 8px;">点击重试</button>
        </div>
      `;
            // 更新超时提示（带重试按钮）
            incomingChatLi.querySelector("p").innerHTML = timeoutMsg;
            incomingRecordLi.querySelector("p").textContent = `请求超时（${window.CONFIG.TIMEOUT_SECONDS || window.CONFIG.TIME_OUT_SECOND/1000}秒）`;
            // 绑定重试事件（复用原有 generateResponse 逻辑）
            incomingChatLi.querySelector(".retry-btn").addEventListener("click", () => {
                incomingChatLi.querySelector("p").textContent = "Thinking...";
                incomingRecordLi.querySelector("p").textContent = "Thinking...";
                generateResponse(incomingChatLi, incomingRecordLi, () => {
                    clearTimeout(timeoutId); // 重试成功后清除超时
                });
            });
        }, window.CONFIG.TIMEOUT_MS || window.CONFIG.TIME_OUT_SECOND);
    };

    // 启动超时定时器
    startTimeout();

    // 发送请求（原有逻辑，补充超时清除）
    generateResponse(incomingChatLi, incomingRecordLi, () => {
        clearTimeout(timeoutId); // 请求成功/失败时清除超时
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
                                    
                                    // 检查content中是否包含航班信息
                                    try {
                                        // 尝试解析content中的JSON数据
                                        if (data.content.includes("flight_number")) {
                                            // 提取JSON部分
                                            const jsonMatch = data.content.match(/\{[\s\S]*?\}/);
                                            if (jsonMatch) {
                                                const flightData = JSON.parse(jsonMatch[0]);
                                                if (flightData.flight_number) {
                                                    // 将单个航班信息转换为数组格式
                                                    renderFlightInfo([flightData]);
                                                }
                                            }
                                        }
                                    } catch (e) {
                                        console.log("No valid flight info in content", e);
                                    }
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

// 渲染航班信息的函数
const renderFlightInfo = (flightInfos) => {
    console.log("Rendering flight info:", flightInfos);
    const flightInfoContainer = document.getElementById('flight-info');
    // 清空容器并初始化基础结构（标题只显示一次）
    flightInfoContainer.innerHTML = `<h3>航班信息</h3>`;

    // 判断参数（数组）是否为空
    if (!flightInfos || flightInfos.length === 0) {
        flightInfoContainer.innerHTML += "<p>没有航班信息可显示</p>";
        return;
    }

    // 循环渲染每一条航班数据
    flightInfos.forEach(flight => {
        // 确保flight是对象
        if (typeof flight === 'string') {
            try {
                flight = JSON.parse(flight);
            } catch (e) {
                console.error("Failed to parse flight info string:", e);
                return;
            }
        }
        
        const flightInfoHTML = `
            <div class="flight-item">
                <p>航班号: ${flight.flight_number || '未知'}</p>
                <p>乘客姓名: ${flight.passengerName || '未填写'}</p>
                <p>出发城市: ${flight.origin || '未知'}</p>
                <p>到达城市: ${flight.destination || '未知'}</p>
                <p>出发时间: ${flight.departure_time || '未知'}</p>
                <p>到达时间: ${flight.arrival_time || '未知'}</p>
                <p>票价: ${flight.price || '未知'}</p>
                <p>座位类型: ${flight.seat_type || '未知'}</p>
                ${flight.message ? `<p class="success-message">状态: ${flight.message}</p>` : ''}
                <hr />
            </div>
        `;
        flightInfoContainer.innerHTML += flightInfoHTML;
    });
    
    // 显示航班信息区域
    flightInfoContainer.style.display = "block";
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