const loginBox = document.getElementById('login-box');
const registerBox = document.getElementById('register-box');
const appView = document.getElementById('app-view');
const authContainer = document.getElementById('auth-container');
let ws = null;
let currentChatTarget = null;
let contextMenuTargetId = null;

// 检查是否已登录
window.onload = function () {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    const userId = localStorage.getItem('user_id');
    if (token) {
        showChat(username, userId);
    }

    // 全局点击关闭右键菜单
    document.addEventListener('click', () => {
        hideContextMenu();
    });
};

function toggleView() {
    loginBox.classList.toggle('hidden');
    registerBox.classList.toggle('hidden');
}

function showChat(username, userId) {
    authContainer.classList.add('hidden');
    appView.classList.remove('hidden');

    if (username && userId) {
        // 更新侧边栏底部的用户信息
        document.getElementById('current-user-name').innerText = username;
        document.getElementById('current-user-id').innerText = `ID: ${userId}`;
        // 如果有头像 URL，也可以在这里设置
        // document.getElementById('current-user-avatar').style.backgroundImage = ...
    }

    // 获取好友列表
    fetchFriends();
    connectWS();
}

// --- 侧边栏功能 ---
function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    const toggleBtn = document.getElementById('sidebar-toggle');
    const icon = toggleBtn.querySelector('i');

    sidebar.classList.toggle('collapsed');

    // 切换图标
    if (sidebar.classList.contains('collapsed')) {
        // 收起状态，图标改为向右箭头（展开）
        icon.classList.remove('fa-chevron-left');
        icon.classList.add('fa-chevron-right');
        toggleBtn.title = "展开侧边栏";
    } else {
        // 展开状态，图标改为向左箭头（收起）
        icon.classList.remove('fa-chevron-right');
        icon.classList.add('fa-chevron-left');
        toggleBtn.title = "收起侧边栏";
    }
}

// --- 好友相关功能 ---

async function fetchFriends() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/v1/relation/list', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        const result = await response.json();
        if (response.ok && result.code === 200) {
            renderFriends(result.data);
        } else {
            console.error("获取好友列表失败:", result.msg);
        }
    } catch (error) {
        console.error("获取好友列表错误:", error);
    }
}

function renderFriends(friends) {
    const list = document.getElementById('friend-list');
    list.innerHTML = '';

    if (!friends || friends.length === 0) {
        list.innerHTML = '<div class="empty-state">暂无好友</div>';
        return;
    }

    friends.forEach(friend => {
        const item = document.createElement('div');
        item.className = 'friend-item';
        if (currentChatTarget == friend.target_id) {
            item.classList.add('active');
        }

        // 左键点击切换聊天
        item.onclick = () => selectFriend(friend.target_id, friend.note, item);

        // 右键点击显示菜单
        item.oncontextmenu = (e) => {
            e.preventDefault();
            showContextMenu(e.pageX, e.pageY, friend.target_id, friend.note);
        };

        const avatar = document.createElement('div');
        avatar.className = 'friend-avatar';
        // 显示名字首字母或默认图标
        const displayName = friend.note || `User ${friend.target_id}`;
        avatar.innerText = displayName.charAt(0).toUpperCase();

        // 在线状态
        const status = document.createElement('div');
        status.className = 'online-status';
        // 使用后端返回的真实在线状态
        if (friend.is_online) {
            status.classList.add('online');
            status.title = "在线";
        } else {
            status.title = "离线";
        }
        avatar.appendChild(status);

        const info = document.createElement('div');
        info.className = 'friend-info';

        const nameSpan = document.createElement('div');
        nameSpan.className = 'friend-name';
        nameSpan.innerText = displayName;

        const idSpan = document.createElement('div');
        idSpan.className = 'friend-id';
        idSpan.innerText = `ID: ${friend.target_id}`;

        info.appendChild(nameSpan);
        info.appendChild(idSpan);

        item.appendChild(avatar);
        item.appendChild(info);
        list.appendChild(item);
    });
}

function selectFriend(targetId, note, domItem) {
    currentChatTarget = targetId;
    const displayName = note || `User ${targetId}`;

    document.getElementById('receiver-id').value = targetId;
    // 更新顶部栏显示的聊天对象名称
    document.getElementById('chat-title').innerText = displayName;

    // 高亮选中状态
    const items = document.querySelectorAll('.friend-item');
    items.forEach(i => i.classList.remove('active'));
    if (domItem) domItem.classList.add('active');

    // 清空消息区域或加载历史消息 (此处仅显示欢迎语)
    const msgList = document.getElementById('message-list');
    msgList.innerHTML = '';

    // 模拟加载历史消息的视觉效果
    const welcome = document.createElement('div');
    welcome.className = 'message-meta';
    welcome.style.textAlign = 'center';
    welcome.style.marginTop = '1rem';
    welcome.innerText = `与 ${displayName} 开始聊天`;
    msgList.appendChild(welcome);
}

// --- 右键菜单功能 ---

function showContextMenu(x, y, targetId, currentNote) {
    contextMenuTargetId = targetId;
    // 存储当前备注以便修改时使用
    document.getElementById('context-menu').dataset.note = currentNote || "";

    const menu = document.getElementById('context-menu');
    menu.style.left = `${x}px`;
    menu.style.top = `${y}px`;
    menu.classList.remove('hidden');
}

function hideContextMenu() {
    const menu = document.getElementById('context-menu');
    menu.classList.add('hidden');
}

function handleModifyRemark() {
    if (!contextMenuTargetId) return;
    const currentNote = document.getElementById('context-menu').dataset.note;
    const newNote = prompt("请输入新的备注:", currentNote);

    if (newNote !== null && newNote !== currentNote) {
        updateFriendNote(contextMenuTargetId, newNote);
    }
    hideContextMenu();
}

function handleDeleteFriend() {
    if (!contextMenuTargetId) return;
    if (confirm(`确定要删除好友 ID:${contextMenuTargetId} 吗？`)) {
        deleteFriend(contextMenuTargetId);
    }
    hideContextMenu();
}

// --- 弹窗功能 ---
function addFriendModal() {
    document.getElementById('add-friend-modal').classList.remove('hidden');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.add('hidden');
}

// --- API 调用 ---

async function updateFriendNote(targetId, newNote) {
    try {
        const token = localStorage.getItem('token');
        const payload = {
            target_id: parseInt(targetId),
            note: newNote
        };

        const response = await fetch('/api/v1/relation/update_note', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        const result = await response.json();
        if (response.ok && result.code === 200) {
            fetchFriends(); // 刷新列表
            // 如果正在聊天的人被修改了备注，更新顶部标题
            if (currentChatTarget == targetId) {
                document.getElementById('chat-title').innerText = newNote || `User ${targetId}`;
            }
        } else {
            alert("更新失败: " + (result.msg || result.error || "未知错误"));
        }
    } catch (error) {
        console.error(error);
        alert("更新出错");
    }
}

async function addFriend() {
    const targetId = document.getElementById('add-friend-id').value;
    const desc = document.getElementById('add-friend-desc').value;

    if (!targetId) {
        alert("请输入好友ID");
        return;
    }

    try {
        const token = localStorage.getItem('token');
        const payload = {
            target_id: parseInt(targetId),
            relation_type: 1, // 1代表好友
            Desc: desc || ""
        };

        const response = await fetch('/api/v1/relation/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        const result = await response.json();
        if (response.ok && result.code === 200) {
            alert("请求发送成功");
            closeModal('add-friend-modal');
            document.getElementById('add-friend-id').value = '';
            document.getElementById('add-friend-desc').value = '';
            fetchFriends(); // 刷新列表
        } else {
            alert("添加失败: " + result.msg);
        }
    } catch (error) {
        console.error(error);
        alert("添加出错");
    }
}

async function deleteFriend(targetId) {
    try {
        const token = localStorage.getItem('token');
        const payload = {
            target_id: parseInt(targetId),
            relation_type: 1
        };

        const response = await fetch('/api/v1/relation/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        const result = await response.json();
        if (response.ok && result.code === 200) {
            fetchFriends(); // 刷新列表
            // 如果删除的是当前聊天对象，清空聊天区域
            if (currentChatTarget == targetId) {
                currentChatTarget = null;
                document.getElementById('chat-title').innerText = "小洋窝"; // 重置标题
                document.getElementById('message-list').innerHTML = '<div class="welcome-message">好友已删除</div>';
            }
        } else {
            alert("删除失败: " + result.msg);
        }
    } catch (error) {
        console.error(error);
        alert("删除出错");
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('user_id');
    if (ws) ws.close();
    location.reload();
}

async function handleAuth(url, data) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        const result = await response.json();

        if (response.ok && result.code === 200) {
            return result;
        } else {
            throw new Error(result.msg || '操作失败');
        }
    } catch (error) {
        alert(error.message);
        throw error;
    }
}

// WebSocket 连接
function connectWS() {
    const token = localStorage.getItem('token');
    if (!token) return;

    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${location.host}/api/v1/chat?token=${token}`;

    if (ws) {
        ws.close();
    }

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log("WebSocket 连接成功");
    };

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        // 如果收到消息，且是当前聊天对象发来的，或者是我发给当前聊天对象的
        // 这里简单处理：只要有消息就追加。更好的做法是判断 sender_id

        let type = 'other';
        let senderName = `User ${msg.sender_id}`;

        // 判断是否是当前聊天窗口
        // 注意：WebSocket 返回的消息可能不包含 sender_id (如果是系统消息)

        if (msg.sender_id) {
            appendMessage(msg.sender_id, msg.content, type);
        }
    };

    ws.onclose = () => {
        console.log("WebSocket 连接断开");
    };

    ws.onerror = (err) => {
        console.error("WS Error:", err);
    };
}

function sendMessage() {
    const receiverID = document.getElementById('receiver-id').value;
    const content = document.getElementById('msg-content').value;

    if (!receiverID || !content) {
        alert("请选择好友或输入ID，并输入内容");
        return;
    }

    const msg = {
        type: "chat",
        receiver_id: parseInt(receiverID),
        content: content
    };

    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(msg));
        appendMessage("Me", content, "self");
        document.getElementById('msg-content').value = '';
    } else {
        alert("连接已断开，正在重连...");
        connectWS();
    }
}

function appendMessage(sender, text, type) {
    const list = document.getElementById('message-list');
    const div = document.createElement('div');
    div.className = `message ${type}`;

    // 如果是对方发的消息，不显示名字在气泡里，而是显示在气泡上方 (meta)
    // 这里为了简单，直接把内容放进去

    div.innerText = text;

    // 如果需要显示发送者名字
    if (type !== 'self') {
        const meta = document.createElement('span');
        meta.className = 'message-meta';
        meta.innerText = sender;
        // 插入到 div 前面？不，HTML结构里我们没有包裹容器。
        // 我们修改一下 DOM 结构：
        // <div wrapper> <span meta> <div bubble> </div>
        // 但为了复用现有 CSS，我们简单点
    }

    list.appendChild(div);
    list.scrollTop = list.scrollHeight;
}

document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const result = await handleAuth('/api/v1/login', { username, password });
        localStorage.setItem('token', result.token);
        localStorage.setItem('username', result.username);
        localStorage.setItem('user_id', result.user_id);
        showChat(result.username, result.user_id);
    } catch (e) {
        console.error(e);
    }
});

document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('reg-username').value;
    const password = document.getElementById('reg-password').value;

    try {
        await handleAuth('/api/v1/register', { username, password });
        alert('注册成功，请登录');
        toggleView();
    } catch (e) {
        console.error(e);
    }
});

// 支持按回车发送
document.getElementById('msg-content').addEventListener('keypress', function (e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});
