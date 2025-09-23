// static/script.js

const BASE_API_URL = '/api';

// --- 功能函数 ---

async function getTodos() {
    try {
        const todosResponse = await fetch(`${BASE_API_URL}/get-all-todos`);
        
        if (!todosResponse.ok) {
            const errorText = await todosResponse.text();
            console.error(`Failed to fetch todos: ${todosResponse.status} ${todosResponse.statusText}. Response:`, errorText);
            try {
                const errorData = JSON.parse(errorText);
                alert(`获取任务失败: ${errorData.error || errorText}`);
            } catch (e) {
                alert(`获取任务失败: ${errorText}`);
            }
            const ul = document.querySelector('ul');
            if (ul) ul.innerHTML = '<p>无法加载任务。</p>';
            return;
        }

        const todosData = await todosResponse.json();
        console.log("Fetched todos data from backend:", todosData);

        const ul = document.querySelector('ul');
        if (!ul) {
            console.error("UL element not found.");
            return;
        }
        ul.innerHTML = '';

        if (!todosData || todosData.length === 0) {
            ul.innerHTML = '<p>没有待办任务，快来创建一个吧！</p>';
            return;
        }

        todosData.forEach(todo => {
            const todoName = todo.name ? todo.name : '无名称';
            const todoDescription = todo.description ? todo.description : '无描述';

            const li = document.createElement('li');
            // ❌ 确保 todo.id 在这里已经被正确赋值，因为是从后端数据中来的
            li.setAttribute('data-todo-id', todo.id); 

            li.innerHTML = `
                <span id="name1-${todo.id}" style="${todo.completed ? 'text-decoration: line-through;' : ''}">${todoName}</span> 
                <span id="description1-${todo.id}" style="${todo.completed ? 'text-decoration: line-through;' : ''}">${todoDescription}</span>
                
                <form class="update-form" id="update-form-${todo.id}" style="display:none;">
                    <input type="text" name="name" placeholder="${todoName}" value="${todo.name || ''}" id="name-${todo.id}">
                    <input type="text" name="description" placeholder="${todoDescription}" value="${todo.description || ''}" id="description-${todo.id}">
                    <button type="submit" class="confirm-update-btn">确认</button> 
                </form>

                <button class="finished-btn">${todo.completed ? '已完成' : '完成'}</button> 
                <button class="display-update-btn">更新</button> 
                <button class="delete-btn">删除</button>
            `;
            
            ul.appendChild(li);
        });
        // ❌ 确保在所有 li 元素都添加到 DOM 之后再设置全局监听器，或者在 DOMContentLoaded 时只设置一次
        // 这里，我们确保在每次 getTodos 更新列表后，事件委托机制仍然有效。
        // 但由于 setupGlobalEventListeners 已经有 DOMContentLoaded 监听，这里不再重复调用。
    } catch (error) {
        console.error("Error fetching todos:", error);
        const ul = document.querySelector('ul');
        if (ul) ul.innerHTML = `<p>加载任务时出现网络或解析错误: ${error.message}</p>`;
    }
}

document.addEventListener('DOMContentLoaded', getTodos);
document.addEventListener('DOMContentLoaded', setupGlobalEventListeners); // 在页面加载时设置一次全局事件监听器

// 修正：确保 setupGlobalEventListeners 是全局可访问的
function setupGlobalEventListeners() {
    const ul = document.querySelector('ul');
    if (!ul) {
        console.error("UL element not found for event delegation.");
        return;
    }

    // ❌ 避免重复绑定事件，先移除所有旧的监听器（如果它已经绑定过）
    // 这对于每次 getTodos() 清空 ul.innerHTML 后重新绑定事件委托尤其重要
    // 但更优雅的作法是确保 ul 不会被 replace 而是更新，这里先用最直接的方法
    // 实际上，只要 setupGlobalEventListeners 只在 DOMContentLoaded 时运行一次，就不会重复绑定。
    // 但是如果 ul.innerHTML 被清空，所有子元素和他们的事件监听就都没了。
    // 所以，这里要确保获取 ID 的逻辑是正确的。

    ul.addEventListener('click', (event) => {
        const target = event.target;
        const li = target.closest('li[data-todo-id]'); // 获取最近的 li 元素
        
        // ❌ 再次确认：如果点击的不是任何带 data-todo-id 的 li 元素内的按钮，直接返回
        if (!li) {
            console.log("Clicked element is not inside a todo LI element.");
            return;
        }

        const todoId = Number(li.getAttribute('data-todo-id')); // 从 li 元素获取 ID
        console.log("Event delegated: Clicked element with todoId:", todoId, "Target class:", target.className); // ❌ 新增日志

        if (isNaN(todoId) || todoId === 0) { // ❌ 增加这里的调试判断
            console.error("Invalid todoId extracted from LI element:", li.getAttribute('data-todo-id'));
            alert("错误：任务 ID 无效，无法执行操作。");
            return;
        }

        if (target.classList.contains('finished-btn')) {
            handleFinished(todoId);
        } else if (target.classList.contains('display-update-btn')) {
            handleDisplay(todoId);
        } else if (target.classList.contains('delete-btn')) {
            deleteTodo(todoId);
        } else if (target.classList.contains('confirm-update-btn')) {
            event.preventDefault(); // 阻止表单默认提交行为
            handleUpdate(todoId);
        }
    });

    // 为表单提交事件也进行委托
    ul.addEventListener('submit', (event) => {
        const target = event.target;
        if (target.classList.contains('update-form')) {
            event.preventDefault(); // 阻止表单默认提交行为
            const li = target.closest('li[data-todo-id]');
            if (!li) {
                console.error("Submit event: LI element not found for form:", target);
                return;
            }
            const todoId = Number(li.getAttribute('data-todo-id'));
            console.log("Event delegated: Form submitted for todoId:", todoId); // ❌ 新增日志
            if (isNaN(todoId) || todoId === 0) { // ❌ 增加这里的调试判断
                console.error("Invalid todoId extracted from form submit:", li.getAttribute('data-todo-id'));
                alert("错误：任务 ID 无效，无法执行操作。");
                return;
            }
            handleUpdate(todoId);
        }
    });
}


async function createTodo() {
    const form = document.querySelector('form');
    const nameInput = form.querySelector('input[name="name"]');
    const descriptionInput = form.querySelector('input[name="description"]');
    
    const name = nameInput.value.trim();
    const description = descriptionInput.value.trim();

    if (!name) {
        alert("任务名称不能为空！");
        return;
    }

    try {
        const response = await fetch(`${BASE_API_URL}/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: name,
                description: description
            })
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Failed to create todo: ${response.status} ${response.statusText}. Response:`, errorText);
            try {
                const errorData = JSON.parse(errorText);
                alert(`创建任务失败: ${errorData.error || errorText}`);
            } catch (e) {
                alert(`创建任务失败: ${errorText}`);
            }
            return;
        }

        form.reset();
        getTodos();
    } catch (error) {
        console.error("Error creating todo:", error);
        alert(`创建任务时出现网络错误: ${error.message}`);
    }
}

function handleDisplay(id) {
    const form = document.querySelector(`#update-form-${id}`);
    if (form) {
        form.style.display = form.style.display === 'none' ? 'block' : 'none'; 
    }
}

async function handleUpdate(id) {
    console.log("handleUpdate received ID (Number):", id, "Is NaN?", isNaN(id));
    
    if (isNaN(id) || id === 0) {
        alert("错误：任务 ID 无效，无法更新。");
        console.error("Invalid ID for handleUpdate:", id);
        return;
    }

    const nameInput = document.querySelector(`#name-${id}`);
    const descriptionInput = document.querySelector(`#description-${id}`);
    
    const name = nameInput ? nameInput.value.trim() : '';
    const description = descriptionInput ? descriptionInput.value.trim() : '';

    if (!name) {
        alert("任务名称不能为空！");
        return;
    }

    try {
        const response = await fetch(`${BASE_API_URL}/update`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: id,
                name: name, 
                description: description,
                completed: false
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Failed to update todo: ${response.status} ${response.statusText}. Response:`, errorText);
            try {
                const errorData = JSON.parse(errorText);
                alert(`更新任务失败: ${errorData.error || errorText}`);
            } catch (e) {
                alert(`更新任务失败: ${errorText}`);
            }
            return;
        }
        
        const form = document.querySelector(`#update-form-${id}`);
        if (form) form.style.display = 'none';
        
        getTodos();
    } catch (error) {
        console.error("Error updating todo:", error);
        alert(`更新任务时出现网络错误: ${error.message}`);
    }
}

async function deleteTodo(id) {
    console.log("deleteTodo received ID (Number):", id, "Is NaN?", isNaN(id));
    
    if (isNaN(id) || id === 0) {
        alert("错误：任务 ID 无效，无法删除。");
        console.error("Invalid ID for deleteTodo:", id);
        return;
    }

    if (!confirm("确定要删除这个任务吗？")) {
        return;
    }
    try {
        const response = await fetch(`${BASE_API_URL}/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({id: id})
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Failed to delete todo: ${response.status} ${response.statusText}. Response:`, errorText);
            try {
                const errorData = JSON.parse(errorText);
                alert(`删除任务失败: ${errorData.error || errorText}`);
            } catch (e) {
                alert(`删除任务失败: ${errorText}`);
            }
            return;
        }
        
        getTodos();
    } catch (error) {
        console.error("Error deleting todo:", error);
        alert(`删除任务时出现网络错误: ${error.message}`);
    }
}

async function handleFinished(id) {
    console.log("handleFinished received ID (Number):", id, "Is NaN?", isNaN(id));

    if (isNaN(id) || id === 0) {
        alert("错误：任务 ID 无效，无法更改状态。");
        console.error("Invalid ID for handleFinished:", id);
        return;
    }

    const nameSpan = document.querySelector(`#name1-${id}`);
    const descriptionSpan = document.querySelector(`#description1-${id}`);
    
    const currentCompleted = nameSpan && nameSpan.style.textDecoration.includes('line-through');
    
    const name = nameSpan ? nameSpan.textContent.trim() : '';
    const description = descriptionSpan ? descriptionSpan.textContent.trim() : '';

    try {
        const response = await fetch(`${BASE_API_URL}/update`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: id,
                name: name, 
                description: description,
                completed: !currentCompleted
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error(`Failed to change status: ${response.status} ${response.statusText}. Response:`, errorText);
            try {
                const errorData = JSON.parse(errorText);
                alert(`标记完成/未完成失败: ${errorData.error || errorText}`);
            } catch (e) {
                alert(`标记完成/未完成失败: ${errorText}`);
            }
            return;
        }
        
        getTodos();
    } catch (error) {
        console.error("Error changing status:", error);
        alert(`标记完成/未完成时出现网络错误: ${error.message}`);
    }
}