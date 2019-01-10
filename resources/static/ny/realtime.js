

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
// var host = "http://192.168.188.240:7777";
var isCountDown = false;
var isEnd = false;
var processBarHight = 30;

function LoadQuestions() {
    $.ajax({
        type: "GET",
        url: "/stream/questions",
        dataType: "json",
        success: function(data){
            ANSWER.questions = data.data.questions;
            ANSWER.curAnswer = data.data.curAnswer;
            DisplayCurQuestion();
        }
    });
}

function DisplayCurQuestion() {
    if (ANSWER.questions.length > 0 && ANSWER.curAnswer < ANSWER.questions.length) {
        // css 重置
        removeAnswerCss();
        var curQuestion = ANSWER.questions[ANSWER.curAnswer];
        $("#AnswerTitle").html((ANSWER.curAnswer+1) + "、" + curQuestion.title);
        $("#btnA").html("A ：" + curQuestion.options[0]);
        $("#btnB").html("B ：" + curQuestion.options[1]);
        $("#btnC").html("C ：" + curQuestion.options[2]);
        $("#btnD").html("D ：" + curQuestion.options[3]);
    } else {
        $("#resetTeam").show();
        $("#resetTeam").click(() => {
            $.ajax({
                type: "GET",
                url: "/stream/reset",
                dataType: "json",
                success: function(data){
                    window.location.reload();
                }
            });
        });
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
}

function answerA(e) {
    var newAnswer = e.data.split(":");
    var now = Date.parse(new Date())/1000;
    if (parseInt(newAnswer[0]) == (ANSWER.curAnswer+1)) {
        if (ANSWER.answerA < now) {
            ANSWER.answerA = now;
        }

        if (ANSWER.answerB == 0 && !isCountDown) {
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

function removeAnswerCss() {
    for (var i=0 ; i < btnMap.length; i++) {
        $("#" + btnMap[i]).removeClass("btn-danger");
        $("#" + btnMap[i]).removeClass("btn-github");
        $("#" + btnMap[i]).removeClass("btn-success");
        $("#" + btnMap[i]).addClass("btn-default");
    }
}

function displayRightAnswer() {
    var rightAid = ANSWER.questions[ANSWER.curAnswer].Answer;
    $("#" + btnMap[rightAid]).removeClass("btn-default");
    $("#" + btnMap[rightAid]).addClass("btn-github");
    var oriHtml = $("#" + btnMap[rightAid]).html();
    $("#" + btnMap[rightAid]).html('<i class="fa fa-github"></i>' + oriHtml);
}

/* jQueryKnob */
$('.knob').knob();

function countDownBar() {
    var time = 0;
    isCountDown = true;
    var interval = setInterval(() => {
        if (time < 15) {
        time++;
        $(".knob").val(time).trigger('change');
        } else{
            countDownFinish();
            clearInterval(interval);
        }
    }, 1000);
}

function countDownFinish() {
    isCountDown = false;
    ANSWER.answerA = 0;
    ANSWER.answerB = 0;
    displayRightAnswer();
    // 显示5秒
    sleep(8).then(() => {
        $(".knob").val(0).trigger('change');
        ANSWER.curAnswer++;
        DisplayCurQuestion();
    });
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms * 1000));
}

/* End jQueryKnob */

function setProcessBarData(data) {
    $("#green-bar").attr("aria-valuenow", data.A);
    $("#green-bar").css("height",parseInt(data.A/processBarHight) + "%");
    if (parseInt(data.A/processBarHight) == 100 && !isEnd) {
        console.log("A END");
        $("#modal-success").modal("show");
        isEnd = true;
    }
    $("#red-bar").attr("aria-valuenow", data.B);
    $("#red-bar").css("height",parseInt(data.B/processBarHight)+ "%");
    if (parseInt(data.B/processBarHight) == 100 && !isEnd) {
        console.log("B END");
        $("#modal-danger").modal("show");
        isEnd = true;
    }
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
