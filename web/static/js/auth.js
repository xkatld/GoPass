// 检查是否已登录
function checkAuth() {
    const token = localStorage.getItem('token');
    if (token) {
        // 如果在登录页面且已有token，跳转到仪表板
        if (window.location.pathname === '/') {
            window.location.href = '/dashboard';
        }
    }
}

// 显示消息
function showMessage(message, type = 'info') {
    const messageDiv = document.getElementById('message');
    const messageIcon = document.getElementById('messageIcon');
    const messageText = document.getElementById('messageText');
    
    messageText.textContent = message;
    messageDiv.className = 'rounded-md p-4';
    
    switch (type) {
        case 'success':
            messageDiv.classList.add('bg-green-50', 'border', 'border-green-200');
            messageIcon.className = 'fas fa-check-circle text-green-400';
            messageText.classList.add('text-green-800');
            break;
        case 'error':
            messageDiv.classList.add('bg-red-50', 'border', 'border-red-200');
            messageIcon.className = 'fas fa-exclamation-circle text-red-400';
            messageText.classList.add('text-red-800');
            break;
        case 'warning':
            messageDiv.classList.add('bg-yellow-50', 'border', 'border-yellow-200');
            messageIcon.className = 'fas fa-exclamation-triangle text-yellow-400';
            messageText.classList.add('text-yellow-800');
            break;
        default:
            messageDiv.classList.add('bg-blue-50', 'border', 'border-blue-200');
            messageIcon.className = 'fas fa-info-circle text-blue-400';
            messageText.classList.add('text-blue-800');
    }
    
    messageDiv.classList.remove('hidden');
    
    // 3秒后自动隐藏
    setTimeout(() => {
        messageDiv.classList.add('hidden');
    }, 3000);
}

// 显示登录表单
function showLoginForm() {
    document.getElementById('loginForm').classList.remove('hidden');
    document.getElementById('registerForm').classList.add('hidden');
}

// 显示注册表单
function showRegisterForm() {
    document.getElementById('loginForm').classList.add('hidden');
    document.getElementById('registerForm').classList.remove('hidden');
}

// 处理登录
async function handleLogin(event) {
    event.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            localStorage.setItem('token', data.data.token);
            localStorage.setItem('user', JSON.stringify(data.data.user));
            showMessage('登录成功！正在跳转...', 'success');
            setTimeout(() => {
                window.location.href = '/dashboard';
            }, 1000);
        } else {
            showMessage(data.message, 'error');
        }
    } catch (error) {
        console.error('Login error:', error);
        showMessage('登录失败，请检查网络连接', 'error');
    }
}

// 处理注册
async function handleRegister(event) {
    event.preventDefault();
    
    const username = document.getElementById('reg-username').value;
    const email = document.getElementById('reg-email').value;
    const password = document.getElementById('reg-password').value;
    
    // 简单的密码强度检查
    if (password.length < 6) {
        showMessage('密码长度至少需要6位', 'warning');
        return;
    }
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            showMessage('注册成功！请登录', 'success');
            showLoginForm();
            // 清空注册表单
            document.getElementById('reg-username').value = '';
            document.getElementById('reg-email').value = '';
            document.getElementById('reg-password').value = '';
        } else {
            showMessage(data.message, 'error');
        }
    } catch (error) {
        console.error('Register error:', error);
        showMessage('注册失败，请检查网络连接', 'error');
    }
}

// 页面加载时检查认证状态
document.addEventListener('DOMContentLoaded', function() {
    checkAuth();
});
