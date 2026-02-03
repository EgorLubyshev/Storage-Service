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
    const username = document.getElementById('loginUsername').value;
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
        "POST", JSON.stringify({"login": username, "password": password}),
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
    const login = document.getElementById('regLogin').value;
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
        "POST", JSON.stringify({"login": login, "password": password, "confirm": confirmPassword}),
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


