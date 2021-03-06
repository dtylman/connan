var os = require('os');

var child;
var fails = 0;
var goBinary = "./connan";

function set_content(html) {
    app.innerHTML = html
    //set focus for autofocus element
    var elem = document.querySelector("input[autofocus]");
    if (elem != null) {
        elem.focus();
    }
}

function show_loader(){
    app.innerHTML += "<div class='loader'><div>";
}

function body_message(msg) {
    set_content('<h1>' + msg + '</h1><div class="loader"></div>');
}

function start_process() {
    if (os.platform().isWindows) {
        goBinary += ".exe";
    }

    body_message("Loading...");

    const spawn = require('child_process').spawn;
    child = spawn(goBinary, { maxBuffer: 1024 * 500 });

    const readline = require('readline');
    const rl = readline.createInterface({
        input: child.stdout
    })

    rl.on('line', (data) => {
        try {
            console.log(`Received: ${data}`);
            if (data.charAt(0) == "$") {
                data = data.substr(1);
                eval(data);
            } else {
                set_content(data);
            }
        } catch (ex) {
            console.log(ex);
        }
    });

    child.stderr.on('data', (data) => {
        console.log(`stderr: ${data}`);
    });

    child.on('close', (code) => {
        alert(`process exited with code ${code}`);
        restart_process();
    });

    child.on('error', (err) => {
        alert('Failed to start child process.');
        restart_process();
    });
}

function restart_process() {
    setTimeout(function () {
        fails++;
        if (fails > 5) {
            close();
        } else {
            start_process();
        }
    }, 5000);
}

function element_as_object(elem) {
    var obj = {
        properties: {}
    }
    for (var j = 0; j < elem.attributes.length; j++) {
        obj.properties[elem.attributes[j].name] = elem.attributes[j].value;
    }
    //overwrite attributes with properties
    if (elem.value != null) {
        obj.properties["value"] = elem.value.toString();
    }
    if (elem.checked != null && elem.checked) {
        obj.properties["checked"] = "true";
    } else {
        delete (obj.properties["checked"]);
    }
    return obj;
}

function element_by_tag_as_array(tag) {
    var items = [];
    var elems = document.getElementsByTagName(tag);
    for (var i = 0; i < elems.length; i++) {
        items.push(element_as_object(elems[i]));
    }
    return items;
}

function fire_event(name, sender) {    
    var msg = {
        name: name,
        sender: element_as_object(sender),
        inputs: element_by_tag_as_array("input").concat(element_by_tag_as_array("select"))
    }
    child.stdin.write(JSON.stringify(msg));
    show_loader();
    console.log(JSON.stringify(msg));
}

function fire_keypressed_event(e, keycode, name, sender) {
    if (e.keyCode === keycode) {
        e.preventDefault();
        fire_event(name, sender);
    }
}

function avoid_reload() {
    if (sessionStorage.getItem("loaded") == "true") {
        alert("go-webkit will fail when page reload. avoid using <form> or submit.");
        close();
    }
    sessionStorage.setItem("loaded", "true");
}

function sidebar_collapse() {
    $('#sidebar').toggleClass('active');
    $('#button-sidebar-collapse').toggleClass('active');
};

function maximize_window() {
    var ngui = require('nw.gui');
    var nwin = ngui.Window.get();
    nwin.show();
    nwin.maximize();
}

function set_progress(progress_id, value, total, label_id, text) {
    pbar = document.getElementById(progress_id)
    if (pbar != null) {
        pbar.innerHTML = value + " / " + total;
        if (total == 0) {
            percentage = 0;
        } else {
            percentage = value * 100 / total;
        }
        pbar.style.width = percentage + "%";
    }
    label = document.getElementById(label_id)
    if (label != null) {
        label.innerHTML = text
    }
}

function attach_scroll_event(id) {
    $(window).scroll(
        function () {
            if ($(window).scrollTop() + $(window).height() > $(document).height() - 1) {
                var input = document.getElementById('scroll-value');
                input.value = $(window).scrollTop();
                document.getElementById(id).click();
            }
        }
    );
}


avoid_reload();
maximize_window();
start_process();