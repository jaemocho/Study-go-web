$(function(){
    if (!window.WebSocket) {
        alert("No WebSocket!")
        return 
    }

    var $chartlog = $('#chat-log')
    var $chatmsg = $('#chat-msg')

    addMessage = function(data) {
        $chartlog.prepend("<div><span>"+ data+"</span></div>");
    }

    connect = function() {
                                    // 현재 윈도우 접속 주소 
        ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onopen = function(e) {
            console.log("onopen", arguments);
        };
        ws.onclose = function(e) {
            console.log("onclose", arguments);
        };
        ws.onmessage = function(e) {
            addMessage(e.data)
            console.log("onmessage", arguments);
        };
    }

    connect();



    var isBlank = function(string) {
        return string == null || string.trim() === "";
    };

    var username;

    while(isBlank(username)) {
        username = prompt("What's your name?")
        if (!isBlank(username)) {
            $('#user-name').html('<b>' + username + '</b>')
        }
    }

    $('#input-form').on('submit', function(e) {

        if (ws.readyState === ws.OPEN) {
            ws.send(JSON.stringify({
                type : "msg", 
                data: $chatmsg.val()
            }));
        }
        $chatmsg.val("")
        $chatmsg.focus()
        // summit button 누를 때 다른 page로 안가려면 false로 
        return false;

    })
}) 