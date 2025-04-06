const chatMessages = document.getElementById('chatMessages');
const userInput = document.getElementById('userInput');
const imageUpload = document.getElementById('imageUpload');
const previewContainer = document.getElementById('previewContainer');
const modelSelect = document.getElementById('model-select');

let selectedModel = modelSelect.value;
let imageFile = null;
let imageBlob = null;

modelSelect.addEventListener("change", (e) => {
    selectedModel = e.target.value;
    const modelChangeMsg = document.createElement("div");
    modelChangeMsg.className = "message ai";
    modelChangeMsg.innerHTML = `
        <div class="avatar"><span class="material-symbols-outlined">smart_toy</span></div>
        <div class="message-content"><p>Model switched to: <strong>${selectedModel}</strong></p></div>
    `;
    chatMessages.appendChild(modelChangeMsg);
    chatMessages.scrollTop = chatMessages.scrollHeight;
});

// ============ Image Handling =============
imageUpload.addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (!file || !file.type.startsWith('image/')) {
        alert('Only image files are supported');
        return;
    }

    imageFile = file;
    await filesToBlob(file);
});

function filesToBlob(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);

        reader.onload = (e) => {
            imageBlob = e.target.result;

            previewContainer.innerHTML = '';
            const preview = document.createElement('div');
            preview.className = 'preview';

            const img = document.createElement('img');
            img.src = imageBlob;

            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'delete-btn';
            deleteBtn.textContent = '×';
            deleteBtn.onclick = clearImage;

            preview.appendChild(img);
            preview.appendChild(deleteBtn);
            previewContainer.appendChild(preview);

            resolve();
        };

        reader.onerror = reject;
    });
}


function clearImage() {
    imageFile = null;
    imageBlob = null;
    previewContainer.innerHTML = '';
    imageUpload.value = '';
}

// ============ Chat Logic =============
function sendMessage() {
    const message = userInput.value.trim();
    if (!message && !imageBlob) return;

    // Display user message
    const userMsg = document.createElement('div');
    userMsg.className = 'message user';
    userMsg.innerHTML = `
        <div class="message-content">
            ${message ? `<p>${message}</p>` : ''}
            ${imageBlob ? `<img src="${imageBlob}" style="width:100px;height:100px;margin-top:6px;border-radius:8px;object-fit:cover;" alt="">` : ''}
        </div>`;
    chatMessages.appendChild(userMsg);
    chatMessages.scrollTop = chatMessages.scrollHeight;

    // Display AI "thinking" message
    const aiMsg = document.createElement("div");
    aiMsg.className = "message ai";
    aiMsg.innerHTML = `
        <div class="avatar"><span class="material-symbols-outlined">smart_toy</span></div>
        <div class="message-content"><p id="streaming-response">Thinking...</p></div>`;
    chatMessages.appendChild(aiMsg);
    chatMessages.scrollTop = chatMessages.scrollHeight;

    let b = imageBlob
    userInput.value = '';
    clearImage();

    // Set timeout control
    const TIMEOUT_MS = 5000; // 5 seconds
    let isTimeout = false;
    const timeoutId = setTimeout(() => {
        isTimeout = true;
        const p = aiMsg.querySelector("#streaming-response");
        if (p) p.textContent = "Request timed out, please try again.";
    }, TIMEOUT_MS);

    generateResponse(message, b, selectedModel, aiMsg, () => {
        if (!isTimeout) clearTimeout(timeoutId);
    });
}

function generateResponse(message, imageBlob, model, containerEl, onFinish) {
    const API_URL = "/api/chat";
    const p = containerEl.querySelector("#streaming-response");
    p.textContent = "";

    let accumulatedResponse = "";

    fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            message,
            bin: imageBlob,
            model
        })
    })
        .then(res => {
            if (!res.ok) throw new Error(`Request failed: ${res.status}`);
            const reader = res.body.getReader();
            const decoder = new TextDecoder();

            function read() {
                return reader.read().then(({ done, value }) => {
                    if (done) {
                        onFinish && onFinish();
                        return;
                    }

                    const chunk = decoder.decode(value);
                    const events = chunk.split('\n\n');

                    events.forEach(event => {
                        if (event.startsWith("event:message")) {
                            const dataLine = event.split('\n').find(line => line.startsWith("data:"));
                            if (dataLine) {
                                try {
                                    const data = JSON.parse(dataLine.replace("data:", "").trim());
                                    accumulatedResponse += data.content;
                                    p.textContent = accumulatedResponse;
                                    chatMessages.scrollTop = chatMessages.scrollHeight;
                                } catch (err) {
                                    console.warn("Parsing failed:", err);
                                }
                            }
                        }
                    });

                    return read();
                });
            }

            return read();
        })
        .catch(err => {
            p.textContent = "An error occurred, please try again later.";
            p.style.color = "red";
            console.error(err);
            onFinish && onFinish();
        });
}

// ============ Input Field Events =============
userInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        sendMessage();
    }
});