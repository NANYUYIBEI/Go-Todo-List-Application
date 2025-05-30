// 获取TODO
async function getTodos() {
    const todos = await fetch('http://localhost:8000/get-all-todos')
    const todosData = await todos.json()
    if (!todosData || todosData.length === 0) {
        const ul = document.querySelector('ul')
        ul.innerHTML = '' // 清空现有内容
        return
    }
    const ul = document.querySelector('ul')
    ul.innerHTML = '' // 清空现有内容
    todosData.forEach(todo => {
        const li = document.createElement('li')
        li.innerHTML = `
            <span id="name1-${todo.id}" style="${todo.completed ? 'text-decoration: line-through;' : ''}">${todo.name}</span> 
            <span id="description1-${todo.id}" style="${todo.completed ? 'text-decoration: line-through;' : ''}">${todo.description}</span>
            <form onsubmit="event.preventDefault(); handleUpdate('${todo.id}')" class="update-form" id="update-form-${todo.id}">
                <input type="text" name="name" placeholder="${todo.name}" id="name-${todo.id}">
                <input type="text" name="description" placeholder="${todo.description}" id="description-${todo.id}">
                <button type="submit">确认</button> 
            </form>
            <button onclick="handleFinished('${todo.id}')">完成</button> 
            <button onclick="handleDisplay('${todo.id}')">更新</button> 
            <button onclick="deleteTodo('${todo.id}')">删除</button>
        `
        ul.appendChild(li)
    })
}

// 获取所有todo
getTodos()

// 创建todo
async function createTodo() {
    const form = document.querySelector('form')
    const formData = new FormData(form)
    const todo = await fetch('http://localhost:8000/create', {
        method: 'POST',
        body: JSON.stringify({
            name: formData.get('name'),
            description: formData.get('description')
        })
    })
    form.reset() // 清空表单
    getTodos()
}

// 显示更新表单
async function handleDisplay(id) {
    const form = document.querySelector(`#update-form-${id}`)
    form.style.display = form.style.display === 'none' ? 'block' : 'none' 
}

// 更新todo
async function handleUpdate(id) {
    const name = document.querySelector(`#name-${id}`).value
    const description = document.querySelector(`#description-${id}`).value
    const todo = await fetch('http://localhost:8000/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: id, 
            name: name, 
            description: description,
            completed: 'false'
        })
    })
    getTodos()
}

// 删除todo
async function deleteTodo(id) {
    const todo = await fetch('http://localhost:8000/delete', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({id: id})
    })
    getTodos()
}

// 完成todo
async function handleFinished(id) {
    const name = document.querySelector(`#name1-${id}`).textContent
    const description = document.querySelector(`#description1-${id}`).textContent

    const todo = await fetch('http://localhost:8000/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: id, 
            name: name, 
            description: description,
            completed: 'true'
        })
    })
    getTodos()
} 