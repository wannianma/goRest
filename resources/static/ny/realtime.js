

function StartRealtime(roomid, timestamp) {
    StartEpoch(timestamp);
    StartSSE();
    LoadQuestions();
    // countDownBar();
}

var ANSWER = {
    "curAnswer": 0,
    "answerA": 0,
    "answerB": 0,
    "questions": {}
};

var btnMap = ["btnA", "btnB", "btnC", "btnD"];
var host = "http://127.0.0.1:7777";

function LoadQuestions() {
    $.ajax({
        type: "GET",
        url: host + "/stream/questions",
        dataType: "json",
        success: function(data){
            ANSWER.questions = data.data;
            DisplayCurQuestion();
        }
    });
}

function DisplayCurQuestion() {
    if (ANSWER.questions.length > 0) {
        var curQuestion = ANSWER.questions[ANSWER.curAnswer];
        $("#AnswerTitle").html((ANSWER.curAnswer+1) + "、" + curQuestion.title);
        $("#btnA").html("A:" + curQuestion.options[0]);
        $("#btnB").html("B:" + curQuestion.options[1]);
        $("#btnC").html("C:" + curQuestion.options[2]);
        $("#btnD").html("D:" + curQuestion.options[3]);
    }
}

function StartEpoch(timestamp) {
    var windowSize = 60;
    var height = 200;
    var defaultData = histogram(windowSize, timestamp);

    window.powerChart = $('#powerChart').epoch({
        type: 'time.line',
        axes: ['bottom', 'left'],
        height: height,
        historySize: 20,
        data: [
            {values: defaultData},
            {values: defaultData}
        ]
    });
}

function StartSSE() {
    if (!window.EventSource) {
        alert("EventSource is not enabled in this browser");
        return;
    }
    var source = new EventSource('/stream/data');
    source.addEventListener('stats', stats, false);
    source.addEventListener('answerA', answerA, false);
    source.addEventListener('answerB', answerB, false);
}

function stats(e) {
    var data = parseJSONStats(e.data);
    powerChart.push(data.power);
    setProcessBarData(data.distance);
    // mallocsChart.push(data.mallocs);
    // messagesChart.push(data.messages);
}

function answerA(e) {
    var newAnswer = e.data.split(":");
    var now = Date.parse(new Date())/1000;
    if (parseInt(newAnswer[0]) == (ANSWER.curAnswer+1)) {
        if (ANSWER.answerA < now) {
            ANSWER.answerA = now;
        }

        if (ANSWER.answerB == 0) {
            countDownBar();
        }
        displayAnswerA(newAnswer[1]);
    }
}

function displayAnswerA(aid) {
    $("#" + btnMap[aid]).removeClass("btn-default");
    $("#" + btnMap[aid]).addClass("btn-success");
}

function answerB(e) {
    var newAnswer = e.data.split(":");
    var now = Date.parse(new Date())/1000;
    if (parseInt(newAnswer[0]) == (ANSWER.curAnswer+1)) {
        if (ANSWER.answerB < now) {
            ANSWER.answerB = now;
        }

        if (ANSWER.answerA == 0) {
            countDownBar();
        }
        displayAnswerB(newAnswer[1]);
    }
}

function displayAnswerB(aid) {
    $("#" + btnMap[aid]).removeClass("btn-default");
    $("#" + btnMap[aid]).addClass("btn-danger");
}


function displayRightAnswer() {
    var rightAid = ANSWER.questions[ANSWER.curAnswer].Answer;
    $("#" + btnMap[rightAid]).removeClass("btn-default");
    $("#" + btnMap[rightAid]).addClass("btn-github");
}

/* jQueryKnob */
$('.knob').knob();

function countDownBar() {
    var time = 0;
    var interval = setInterval(() => {
        if (time < 10 * 2) {
        time++;
        $(".knob").val(time).trigger('change');
        } else{
            countDownFinish();
            clearInterval(interval);
        }
    }, 1000);
}

function countDownFinish() {
    ANSWER.answerA = 0;
    ANSWER.answerB = 0;
    
    $(".knob").val(0).trigger('change');
}

/* End jQueryKnob */

function setProcessBarData(data) {
    $("#green-bar").attr("aria-valuenow", data.A);
    $("#green-bar").css("height",parseInt(data.A/10) + "%");

    $("#red-bar").attr("aria-valuenow", data.B);
    $("#red-bar").css("height",parseInt(data.B/10)+ "%");
}

function parseJSONStats(e) {
    var data = jQuery.parseJSON(e);
    var timestamp = data.timestamp;

    var power = [
        {time: timestamp, y: data.powerA},
        {time: timestamp, y: data.powerB}
    ];

    var distance = {
        A: data.distanceA,
        B: data.distanceB
    }

    return {
        power: power,
        distance: distance
        // messages: messages
    }
}

function histogram(windowSize, timestamp) {
    var entries = new Array(windowSize);
    for(var i = 0; i < windowSize; i++) {
        entries[i] = {time: (timestamp-windowSize+i-1), y:0};
    }
    return entries;
}

var entityMap = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': '&quot;',
    "'": '&#39;',
    "/": '&#x2F;'
};

function rowStyle(nick) {
    var classes = ['active', 'success', 'info', 'warning', 'danger'];
    var index = hashCode(nick)%5;
    return classes[index];
}

function hashCode(s){
  return Math.abs(s.split("").reduce(function(a,b){a=((a<<5)-a)+b.charCodeAt(0);return a&a},0));             
}

function escapeHtml(string) {
    return String(string).replace(/[&<>"'\/]/g, function (s) {
      return entityMap[s];
    });
}

window.StartRealtime = StartRealtime
