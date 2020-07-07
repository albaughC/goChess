//https://www.learnwithjason.dev/blog/get-form-values-as-json/
const formToJson = elements => [].reduce.call(elements, (data, element) => {
    if (isValidElement(element)) {
        data[element.name] = element.value;
    }
    return data;
}, {});

const handleFormSubmit = event => {
    event.preventDefault();
    const data = formToJson(event.target.elements);
    fetch('http://127.0.0.1:8000/api/login', {
        method:'POST',
        headers: {
            'Accept': 'application/json, text/plain, */*',
            'Content-Type': 'application/json'
            },
        body: JSON.stringify(data)
    })
    .then((res) => res.text())
    .then((data) => console.log(data))
};

const isValidElement = element => {
    return element.name && element.value;
}
            
const loginForm = document.querySelector('.login-form');
loginForm.addEventListener('submit', handleFormSubmit);

const registerForm = document.querySelector('.register-form');
registerForm.addEventListener('submit', handleFormSubmit);

document.querySelector('#log-link').addEventListener('click', event => {
    document.querySelector('#login-container').style.display = "none";
    document.querySelector('#register-container').style.display = "flex";
})
document.querySelector('#reg-link').addEventListener('click', event => {
    document.querySelector('#login-container').style.display = "flex";
    document.querySelector('#register-container').style.display= "none";
})
