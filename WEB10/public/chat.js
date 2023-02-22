$(function(){
    if (!window.EventSource) {
        alert("No EventSource!")
        return 
    }

    var $chartlog = $('#chat-log')
    var $chatmsg = $('#chat-msg')

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
        $.post('/messages', {
            msg : $chatmsg.val(),
            name: username
        });
        $chatmsg.val("")
        $chatmsg.focus()
        // summit button 누를 때 다른 page로 안가려면 false로 
        return false;

    })

    var addMessage = function(data) {
        var text = "";
        if (!isBlank(data.name)){
            text = '<strong> ' + data.name+ ':</strong>';
        }
        text += data.msg;
        $chartlog.prepend('<div><span>' + text + '</span></div>');
    }

    // addMessage({
    //     msg: 'hello',
    //     name: 'aaa'
    // })

    // addMessage({
    //     msg: 'hello2',
    // })

    // event source url 지정 
    var es = new EventSource('/stream')
    es.onopen = function(e) {
        $.post('users/', {
            name: username
        })
    }

    // eventsource 에 message가 왔을 대 동작
    es.onmessage = function(e) {
        var msg = JSON.parse(e.data);
        addMessage(msg);
    }

    // 화면이 닫히기 직전에 수행 
    window.onbeforeunload = function() {
        $.ajax({
            url: "/users>username=" + username,
            type: "DELETE"
        });
        es.close();
    }

}) 