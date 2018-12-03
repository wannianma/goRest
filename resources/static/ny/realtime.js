

function StartRealtime(roomid, timestamp) {
    StartEpoch(timestamp);
    StartSSE();
    countDownBar();
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
}

function stats(e) {
    var data = parseJSONStats(e.data);
    powerChart.push(data.power);
    setProcessBarData(data.distance);
    // mallocsChart.push(data.mallocs);
    // messagesChart.push(data.messages);
}

/* jQueryKnob */
$('.knob').knob();

function countDownBar() {
    var time = 0;
    var interval = setInterval(() => {
        if (time < 10) {
        time++;
        $(".knob").val(time).trigger('change');
        } else{
            countDownFinish();
            clearInterval(interval);
        }
    }, 1000);
}

function countDownFinish() {
    alert("10 second finish");
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
