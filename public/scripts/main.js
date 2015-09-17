var StatSAMP = function() {
  var serversHistoryPeriod = 'day',
      playersHistoryPeriod = 'day',
      serversHistoryChart = null,
      playersHistoryChart = null;

  // Updating data
  var updateGeneral = function() {
    $.get("/api/general", function(data) {
      $("#servers-total-count").text(data.servers_total);
      $("#servers-online-count").text(data.servers_online);
      $("#servers-offline-count").text(data.servers_offline);
      renderServersPie([
        {
            value: data.servers_online,
            color:"rgba(33, 150, 243, .75)",
            highlight: "rgb(33, 150, 243)",
            label: "Online"
        },
        {
            value: data.servers_offline,
            color: "rgba(255,171,64, .75)",
            highlight: "rgb(255,171,64)",
            label: "Offline"
        }
      ]);
      
      $("#slots-total-count").text(data.slots_total);
      $("#slots-used-count").text(data.slots_used);
      $("#slots-unused-count").text(data.slots_unused);
      renderPlayersPie([
        {
            value: data.slots_used,
            color:"rgba(33, 150, 243, .75)",
            highlight: "rgb(33, 150, 243)",
            label: "Used"
        },
        {
            value: data.slots_unused,
            color: "rgba(255,171,64, .75)",
            highlight: "rgb(255,171,64)",
            label: "Unused"
        }
      ]);
    });
  };

  var updateServersHistory = function() {
    $('#servers-chart-period-tooltip').html('Period: <b>' + serversHistoryPeriod + '</b>');
    $.get("/api/history/servers/" + serversHistoryPeriod, function(data) {
      var chartData = {
        responsive: true,
        maintainAspectRatio: true,
        labels: [],
        datasets: [
          {
            fillColor: "rgba(33, 150, 243, .25)",
            strokeColor: "rgb(33, 150, 243)",
            pointColor: "rgb(33, 150, 243)",
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgb(33, 150, 243)",
            data: []
          }
        ]
      };
      $.each(data, function(index, value) {
        chartData.labels.push(formatTime(value.date));
        chartData.datasets[0].data.push(value.total);
      });
      renderServersHistoryChart(chartData);
    });
  };

  var updatePlayersHistory = function() {
    $('#players-chart-period-tooltip').html('Period: <b>' + playersHistoryPeriod + '</b>');
    $.get("/api/history/players/" + playersHistoryPeriod, function(data) {
      var chartData = {
        responsive: true,
        maintainAspectRatio: true,
        labels: [],
        datasets: [
          {
            fillColor: "rgba(33, 150, 243, .25)",
            strokeColor: "rgb(33, 150, 243)",
            pointColor: "rgb(33, 150, 243)",
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgb(33, 150, 243)",
            data: []
          }
        ]
      };
      $.each(data, function(index, value) {
        chartData.labels.push(formatTime(value.date));
        chartData.datasets[0].data.push(value.total);
      });
      renderPlayersHistoryChart(chartData);
    });
  };

  // Render charts

  var renderServersPie = function(data) {
    var context = document.getElementById("servers-pie").getContext("2d");
    window.serversPie = new Chart(context).Doughnut(data, { percentageInnerCutout : 75 });
  };
  
  var renderPlayersPie = function(data) {
    var context = document.getElementById("players-pie").getContext("2d");
    window.playersPie = new Chart(context).Doughnut(data, { percentageInnerCutout : 75 });
  };

  var renderServersHistoryChart = function(data) {
    var context = document.getElementById("servers-history-chart").getContext("2d");
    if(serversHistoryChart != null) {
      serversHistoryChart.destroy();
      serversHistoryChart = null;
    }
    serversHistoryChart = new Chart(context).Line(data, {
      responsive: true,
      tooltipTemplate: "<%= value %> servers"
    });
  };

  var renderPlayersHistoryChart = function(data) {
    var context = document.getElementById("players-history-chart").getContext("2d");
    if(playersHistoryChart != null) {
      playersHistoryChart.destroy();
      playersHistoryChart = null;
    }
    playersHistoryChart = new Chart(context).Line(data, {
      responsive: true,
      tooltipTemplate: "<%= value %> servers"
    });
  };
  
  function formatTime(datestring) {
    var date = new Date(datestring);
    return date.getHours() + ':' + date.getMinutes();
  };
  
  function periodValueToString(value)
  {
    switch(value) {
      case '1':
        return "day";
      case '2':
        return "week";
      case '3':
        return "month";
      case '4':
        return "year";
      default:
        return "day";
    }
  }

  // Events
  $(document).ready(function() {
    updateServersHistory();
    updatePlayersHistory();
    updateGeneral();
  });
  
  $('#servers-chart-period').change(function() {
    var value = $(this).val();
    serversHistoryPeriod = periodValueToString(value);
    updateServersHistory();
  });
  
  $('#players-chart-period').change(function() {
    var value = $(this).val();
    playersHistoryPeriod = periodValueToString(value);
    updatePlayersHistory();
  });
};

StatSAMP();
