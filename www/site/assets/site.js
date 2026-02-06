let App

function onStart() {
    App = {UserToken: null}
    Ajax(
        'http://localhost:8080/api/v1/user/login',
        "GET", null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
        },
        function (x) {
            alert("Cancel: " + x.status);
        },
    )
}

function onLoginClick() {
    alert("LoginClick");
}

function GetLoginForm() {
    onStart()
}

function GetRegisterForm() {


    Ajax(
        'http://localhost:8080/api/v1/user/register',
        "GET", null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
        },
        function (x) {
            alert("Cancel: " + x.status);

        },
    )

}


function LoginClick() {
    const username = document.getElementById('loginName').value;
    const password = document.getElementById('loginPassword').value;

    if (!username || !password) {
        alert('Пожалуйста, заполните все поля');
        return;
    }

    // Здесь будет отправка данных на сервер
    console.log('Вход:', {
        username,
        password
    });

    Ajax(
        'http://localhost:8080/api/v1/user/login',
        "POST", JSON.stringify({"name": username, "password": password}),
        function (x) {
            App.UserToken = x.getResponseHeader("x-user-token");
         
            getLobby();
        },
        function (x) {
            console.log(x.responseText)
            onStart();
        },
    )
}


function RegClick() {
    const login = document.getElementById('regName').value;
    const password = document.getElementById('regPassword').value;
    const confirmPassword = document.getElementById('regConfirmPassword').value;

    if (!login || !password || !confirmPassword) {
        alert('Пожалуйста, заполните все поля');
        return;
    }

    // Здесь будет отправка данных на сервер
    console.log('Вход:', {
        login,
        password,
        confirmPassword
    });

    Ajax(
        'http://localhost:8080/api/v1/user/register',
        "POST", JSON.stringify({"name": login, "password": password, "confirm": confirmPassword}),
        function (x) {
            console.log("Регистрация завершена");

            App.UserToken = x.getResponseHeader("x-user-token");
            getLobby();
        },
        function (x) {
            alert(x.responseText)
            onStart();
        },
    )
}

function getLobby() {
    Ajax(
        'http://localhost:8080/api/v1/lobby',
        "GET",null,
        function (x) {
           
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
            if (App.UserToken) {
                console.log("TOKEN:", App.UserToken);
            }
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
}

function SendLobbyData() {
    const input = document.getElementById('lobbyData');
    const status = document.getElementById('lobbyDataStatus');
    const data = input ? input.value.trim() : '';

    if (!data) {
        alert('Необходимо ввести данные для отправки');
        return;
    }

    Ajax(
        'http://localhost:8080/api/v1/lobby/data',
        "POST", JSON.stringify({"data": data}),
        function (x) {
            const resp = JSON.parse(x.responseText);
            if (status) {
                status.innerText = "Отправлено: " + resp.modified + " (len=" + resp.length + ")";
            }
        },
        function (x) {
            if (status) {
                status.innerText = "Ошибка: " + x.responseText;
            }
        },
    )
}

function UploadLobbyFiles() {
    const input = document.getElementById('lobbyFiles');
    const status = document.getElementById('lobbyFilesStatus');

    if (!input || !input.files || input.files.length === 0) {
        alert(' Пожалуйста, выберите файлы для загрузки');
        return;
    }

    const form = new FormData();
    for (let i = 0; i < input.files.length; i++) {
        form.append('files', input.files[i]);
    }

    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status >= 200 && xhr.status < 300) {
                const resp = JSON.parse(xhr.responseText);
                if (status) {
                    const names = resp.files.map(f => f.name + " (" + f.size + " bytes)").join(", ");
                    status.innerText = "Загружено: " + names;
                }
            } else {
                if (status) {
                    status.innerText = "Ошибка: " + xhr.responseText;
                }
            }
        }
    };

    xhr.open("POST", "http://localhost:8080/api/v1/lobby/files", false);
    if (App.UserToken != null) {
        xhr.setRequestHeader("Authorization", "Bearer " + App.UserToken);
    }
    xhr.send(form);
}
 function  CreateGame(){
    Ajax(
        'http://localhost:8080/api/v1/games/create',
        "POST",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            
            CreateGameHtml();
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
 }
 function CreateGameHtml(){
    Ajax(
        'http://localhost:8080/api/v1/games/create/html',
        "GET",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
           
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
 }
 function CheckGame(id){
    Ajax(
        'http://localhost:8080/api/v1/games/'+id,
        "GET",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
            MoveHtml(id);
        },
        function (x) {
            alert("2й игрок не появился");
          
        },
    )
 }
 function ConnectToGame() {
    Ajax(
        'http://localhost:8080/api/v1/games/list',
        "GET",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
}
 function startGame (id) {
    Ajax(
        'http://localhost:8080/api/v1/games/'+id+'/connect',
        "POST",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
            MoveHtml(id);
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
}
 
function MoveHtml(id){
    Ajax(
        'http://localhost:8080/api/v1/games/'+id+'/move/html',
        "GET",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
           
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
}

function Move(index, id){
    Ajax(
        'http://localhost:8080/api/v1/games/'+id+'/move',
        "POST",JSON.stringify({"index":parseInt(index)}),
        function (x) {
 
            MoveHtml(id);
        },
        function (x) {
            alert("Не твой ход");
          
        },
    )
}
function Exit(id){
    Ajax(
        'http://localhost:8080/api/v1/games/'+id+'/exit/html',
        "GET",null,
        function (x) {
            const resp = JSON.parse(x.responseText);
            document.getElementById('output').innerHTML = resp.html;
           
        },
        function (x) {
            alert("Cancel: " + x.status);
          
        },
    )
}

onStart();


