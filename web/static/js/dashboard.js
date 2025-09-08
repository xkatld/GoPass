let passwords = [];
let allPasswords = []; // 保存所有密码用于搜索
let categories = [];
let currentEditingId = null;

// 检查认证状态
function checkAuth() {
    const token = localStorage.getItem('token');
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    
    if (!token) {
        window.location.href = '/';
        return false;
    }
    
    // 显示用户名
    document.getElementById('username').textContent = user.username || '用户';
    return true;
}

// 获取认证头
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

// 显示Toast消息
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    const toastIcon = document.getElementById('toastIcon');
    const toastMessage = document.getElementById('toastMessage');
    
    toastMessage.textContent = message;
    
    switch (type) {
        case 'success':
            toastIcon.className = 'fas fa-check-circle text-green-500';
            break;
        case 'error':
            toastIcon.className = 'fas fa-exclamation-circle text-red-500';
            break;
        case 'warning':
            toastIcon.className = 'fas fa-exclamation-triangle text-yellow-500';
            break;
        default:
            toastIcon.className = 'fas fa-info-circle text-blue-500';
    }
    
    toast.classList.remove('hidden');
    
    setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// 加载密码列表
async function loadPasswords() {
    try {
        const response = await fetch('/api/passwords', {
            headers: getAuthHeaders()
        });
        
        const data = await response.json();
        
        if (data.success) {
            allPasswords = data.data || [];
            passwords = [...allPasswords];
            renderPasswordList();
            loadCategories();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        console.error('Load passwords error:', error);
        showToast('加载密码列表失败', 'error');
    }
}

// 渲染密码列表
function renderPasswordList() {
    const passwordList = document.getElementById('passwordList');
    const emptyState = document.getElementById('emptyState');
    
    if (passwords.length === 0) {
        passwordList.innerHTML = '';
        emptyState.classList.remove('hidden');
        return;
    }
    
    emptyState.classList.add('hidden');
    
    passwordList.innerHTML = passwords.map(password => {
        const websiteUrl = password.website && !password.website.startsWith('http') ?
            `https://${password.website}` : password.website;
        const domain = password.website ? new URL(websiteUrl || 'https://example.com').hostname : '';
        const faviconUrl = domain ? `https://www.google.com/s2/favicons?domain=${domain}&sz=32` : '';

        return `
        <div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow duration-200">
            <div class="flex items-start justify-between mb-4">
                <div class="flex items-center space-x-3">
                    <div class="flex-shrink-0">
                        ${faviconUrl ?
                            `<img src="${faviconUrl}" alt="网站图标" class="w-8 h-8 rounded" onerror="this.style.display='none'; this.nextElementSibling.style.display='block';">
                             <div class="w-8 h-8 bg-gray-200 rounded flex items-center justify-center" style="display:none;">
                                <i class="fas fa-globe text-gray-400 text-sm"></i>
                             </div>` :
                            `<div class="w-8 h-8 bg-gray-200 rounded flex items-center justify-center">
                                <i class="fas fa-globe text-gray-400 text-sm"></i>
                             </div>`
                        }
                    </div>
                    <div class="flex-1 min-w-0">
                        <h3 class="text-lg font-medium text-gray-900 truncate">${password.title}</h3>
                        <p class="text-sm text-gray-500 truncate">${password.username || '无用户名'}</p>
                    </div>
                </div>
                <div class="flex items-center space-x-2">
                    <button onclick="copyPassword('${password.password}')"
                            class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
                            title="复制密码">
                        <i class="fas fa-copy"></i>
                    </button>
                    <button onclick="editPassword(${password.id})"
                            class="p-2 text-indigo-600 hover:text-indigo-800 hover:bg-indigo-50 rounded-full transition-colors"
                            title="编辑">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button onclick="deletePassword(${password.id})"
                            class="p-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-full transition-colors"
                            title="删除">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>

            <div class="space-y-2">
                ${password.website ? `
                <div class="flex items-center text-sm text-gray-600">
                    <i class="fas fa-link w-4 mr-2"></i>
                    <a href="${websiteUrl}" target="_blank" class="text-blue-600 hover:text-blue-800 truncate">
                        ${password.website}
                    </a>
                </div>` : ''}

                ${password.category ? `
                <div class="flex items-center text-sm text-gray-600">
                    <i class="fas fa-tag w-4 mr-2"></i>
                    <span class="px-2 py-1 bg-gray-100 text-gray-800 rounded-full text-xs">
                        ${password.category}
                    </span>
                </div>` : ''}

                ${password.notes ? `
                <div class="flex items-start text-sm text-gray-600">
                    <i class="fas fa-sticky-note w-4 mr-2 mt-0.5"></i>
                    <span class="truncate">${password.notes}</span>
                </div>` : ''}
            </div>

            <div class="mt-4 pt-4 border-t border-gray-100">
                <div class="flex items-center justify-between text-xs text-gray-400">
                    <span>创建时间: ${new Date(password.created_at).toLocaleDateString()}</span>
                    <span>ID: ${password.id}</span>
                </div>
            </div>
        </div>
        `;
    }).join('');
}

// 复制密码到剪贴板
async function copyPassword(password) {
    try {
        await navigator.clipboard.writeText(password);
        showToast('密码已复制到剪贴板', 'success');
    } catch (error) {
        console.error('Copy error:', error);
        showToast('复制失败', 'error');
    }
}

// 显示添加密码模态框
function showAddPasswordModal() {
    currentEditingId = null;
    document.getElementById('modalTitle').textContent = '添加密码';
    document.getElementById('passwordForm').reset();
    document.getElementById('passwordId').value = '';
    document.getElementById('passwordModal').classList.remove('hidden');
}

// 编辑密码
function editPassword(id) {
    const password = passwords.find(p => p.id === id);
    if (!password) return;
    
    currentEditingId = id;
    document.getElementById('modalTitle').textContent = '编辑密码';
    document.getElementById('passwordId').value = id;
    document.getElementById('title').value = password.title;
    document.getElementById('website').value = password.website || '';
    document.getElementById('passwordUsername').value = password.username || '';
    document.getElementById('passwordField').value = password.password;
    document.getElementById('category').value = password.category || '';
    document.getElementById('notes').value = password.notes || '';
    
    document.getElementById('passwordModal').classList.remove('hidden');
}

// 关闭密码模态框
function closePasswordModal() {
    document.getElementById('passwordModal').classList.add('hidden');
    currentEditingId = null;
}

// 处理密码表单提交
async function handlePasswordSubmit(event) {
    event.preventDefault();
    
    const formData = {
        title: document.getElementById('title').value,
        website: document.getElementById('website').value,
        username: document.getElementById('passwordUsername').value,
        password: document.getElementById('passwordField').value,
        category: document.getElementById('category').value,
        notes: document.getElementById('notes').value
    };
    
    try {
        const url = currentEditingId ? `/api/passwords/${currentEditingId}` : '/api/passwords';
        const method = currentEditingId ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: getAuthHeaders(),
            body: JSON.stringify(formData)
        });
        
        const data = await response.json();
        
        if (data.success) {
            showToast(currentEditingId ? '密码更新成功' : '密码添加成功', 'success');
            closePasswordModal();
            await loadPasswords();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        console.error('Save password error:', error);
        showToast('保存失败', 'error');
    }
}

// 删除密码
async function deletePassword(id) {
    if (!confirm('确定要删除这个密码条目吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/passwords/${id}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });
        
        const data = await response.json();
        
        if (data.success) {
            showToast('密码删除成功', 'success');
            await loadPasswords();
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        console.error('Delete password error:', error);
        showToast('删除失败', 'error');
    }
}

// 切换密码可见性
function togglePasswordVisibility(fieldId) {
    const field = document.getElementById(fieldId);
    const icon = field.nextElementSibling.querySelector('i');
    
    if (field.type === 'password') {
        field.type = 'text';
        icon.className = 'fas fa-eye-slash text-gray-400';
    } else {
        field.type = 'password';
        icon.className = 'fas fa-eye text-gray-400';
    }
}

// 显示生成密码模态框
function showGeneratePasswordModal() {
    document.getElementById('generatePasswordModal').classList.remove('hidden');
    generateNewPassword();
}

// 关闭生成密码模态框
function closeGeneratePasswordModal() {
    document.getElementById('generatePasswordModal').classList.add('hidden');
}

// 更新长度显示
function updateLengthDisplay() {
    const length = document.getElementById('passwordLength').value;
    document.getElementById('lengthDisplay').textContent = length;
}

// 生成新密码
async function generateNewPassword() {
    const length = parseInt(document.getElementById('passwordLength').value);
    const includeUpper = document.getElementById('includeUpper').checked;
    const includeLower = document.getElementById('includeLower').checked;
    const includeNumbers = document.getElementById('includeNumbers').checked;
    const includeSymbols = document.getElementById('includeSymbols').checked;
    
    try {
        const response = await fetch('/api/generate-password', {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                length,
                include_upper: includeUpper,
                include_lower: includeLower,
                include_numbers: includeNumbers,
                include_symbols: includeSymbols
            })
        });
        
        const data = await response.json();
        
        if (data.success) {
            document.getElementById('generatedPassword').value = data.data.password;
            displayPasswordStrength(data.data.strength);
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        console.error('Generate password error:', error);
        showToast('生成密码失败', 'error');
    }
}

// 显示密码强度
function displayPasswordStrength(strength) {
    const strengthDiv = document.getElementById('passwordStrength');
    const score = strength.score;
    const strengthText = strength.strength;
    
    let colorClass = '';
    let strengthLabel = '';
    
    switch (strengthText) {
        case 'weak':
            colorClass = 'text-red-600';
            strengthLabel = '弱';
            break;
        case 'medium':
            colorClass = 'text-yellow-600';
            strengthLabel = '中等';
            break;
        case 'strong':
            colorClass = 'text-green-600';
            strengthLabel = '强';
            break;
    }
    
    strengthDiv.innerHTML = `
        <div class="flex items-center justify-between">
            <span class="${colorClass} font-medium">强度: ${strengthLabel}</span>
            <span class="text-gray-500">评分: ${score}/6</span>
        </div>
        <div class="mt-1 bg-gray-200 rounded-full h-2">
            <div class="h-2 rounded-full ${colorClass.replace('text-', 'bg-')}" style="width: ${(score/6)*100}%"></div>
        </div>
    `;
}

// 复制到剪贴板
async function copyToClipboard(elementId) {
    const element = document.getElementById(elementId);
    try {
        await navigator.clipboard.writeText(element.value);
        showToast('已复制到剪贴板', 'success');
    } catch (error) {
        console.error('Copy error:', error);
        showToast('复制失败', 'error');
    }
}

// 导出数据
async function exportData() {
    try {
        const response = await fetch('/api/export', {
            headers: getAuthHeaders()
        });

        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `gopass_export_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.csv`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
            showToast('数据导出成功', 'success');
        } else {
            showToast('导出失败', 'error');
        }
    } catch (error) {
        console.error('Export error:', error);
        showToast('导出失败', 'error');
    }
}

// 显示导入模态框
function showImportModal() {
    document.getElementById('importModal').classList.remove('hidden');
    // 关闭导入导出菜单
    document.getElementById('importExportMenu').classList.add('hidden');
}

// 关闭导入模态框
function closeImportModal() {
    document.getElementById('importModal').classList.add('hidden');
    document.getElementById('importForm').reset();
}

// 处理导入
async function handleImport(event) {
    event.preventDefault();

    const fileInput = document.getElementById('importFile');
    const file = fileInput.files[0];

    if (!file) {
        showToast('请选择文件', 'warning');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch('/api/import', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: formData
        });

        const data = await response.json();

        if (data.success) {
            showToast(data.message, 'success');
            closeImportModal();
            await loadPasswords(); // 重新加载密码列表
        } else {
            showToast(data.message, 'error');
        }
    } catch (error) {
        console.error('Import error:', error);
        showToast('导入失败', 'error');
    }
}

// 切换导入导出菜单
function toggleImportExportMenu() {
    const menu = document.getElementById('importExportMenu');
    menu.classList.toggle('hidden');
}

// 加载分类
async function loadCategories() {
    try {
        const response = await fetch('/api/categories', {
            headers: getAuthHeaders()
        });

        const data = await response.json();

        if (data.success) {
            categories = data.data || [];
            renderCategoryFilter();
        }
    } catch (error) {
        console.error('Load categories error:', error);
    }
}

// 渲染分类筛选器
function renderCategoryFilter() {
    const categoryFilter = document.getElementById('categoryFilter');
    const uniqueCategories = [...new Set(allPasswords.map(p => p.category).filter(c => c))];

    categoryFilter.innerHTML = '<option value="">所有分类</option>';
    uniqueCategories.forEach(category => {
        categoryFilter.innerHTML += `<option value="${category}">${category}</option>`;
    });
}

// 搜索和筛选
function filterPasswords() {
    const searchTerm = document.getElementById('searchInput').value.toLowerCase();
    const selectedCategory = document.getElementById('categoryFilter').value;

    passwords = allPasswords.filter(password => {
        const matchesSearch = !searchTerm ||
            password.title.toLowerCase().includes(searchTerm) ||
            (password.website && password.website.toLowerCase().includes(searchTerm)) ||
            (password.username && password.username.toLowerCase().includes(searchTerm));

        const matchesCategory = !selectedCategory || password.category === selectedCategory;

        return matchesSearch && matchesCategory;
    });

    renderPasswordList();
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    if (checkAuth()) {
        loadPasswords();
    }

    // 添加搜索和筛选事件监听器
    const searchInput = document.getElementById('searchInput');
    const categoryFilter = document.getElementById('categoryFilter');

    if (searchInput) {
        searchInput.addEventListener('input', filterPasswords);
    }

    if (categoryFilter) {
        categoryFilter.addEventListener('change', filterPasswords);
    }

    // 点击外部关闭菜单
    document.addEventListener('click', function(event) {
        const menu = document.getElementById('importExportMenu');
        const button = event.target.closest('button');
        if (!button || !button.onclick || button.onclick.toString().indexOf('toggleImportExportMenu') === -1) {
            menu.classList.add('hidden');
        }
    });
});
