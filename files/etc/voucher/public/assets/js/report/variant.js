$(document).ready(function() {
  getVariant();
});

function getVariant() {
    console.log("Get Variant Data");
    $.ajax({
        url: '/v1/report/variant',
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;
          initFlotBar(result);
        },
        error: function (data) {
          alert("Variant Not Found.");
        }
    });
}

function getVoucher(id) {
    console.log("Get Variant Data");
    $.ajax({
        url: '/v1/report/voucher/variant?id='+id,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;

          console.log(result);
          initFlotBar(result);
        },
        error: function (data) {
          alert("Variant Not Found.");
        }
    });
}

function initFlotBar(data){
  // Flot bar chart
  // ------------------
  var barChartOptions = {
      series: {
          bars: {
              show: true,
              fill: 1,
              barWidth: 0.2,
              lineWidth: 0,
              align: 'center'
          },
          highlightColor: 'rgba(255,255,255,0.2)'
      },
      grid: {
          hoverable: true,
          clickable: true,
          borderWidth: 0,
          color: '#8394a9'
      },
      tooltip: true,
      tooltipOpts: {
          content: function getTooltip(label, x, y) {
              return 'Program for ' + x + ' are ' + y;
          }
      },
      xaxis: {
          tickColor: 'transparent',
          mode: 'categories',
          font: {
              color: Colors.byName('blueGrey-200')
          }
      },
      yaxis: {
          tickColor: 'transparent',
          font: {
              color: Colors.byName('blueGrey-200')
          }
      },
      legend: {
          show: false
      },
      shadowSize: 0
  };

  var barChartData = [{
      'label': 'New',
      bars: {
          order: 0,
          fillColor: Colors.byName('primary')
      },
      data: [
          ['Jan', 20],
          ['Feb', 15],
          ['Mar', 25],
          ['Apr', 30],
          ['May', 40],
          ['Jun', 35]
      ]
  }, {
      'label': 'Recurrent',
      bars: {
          order: 1,
          fillColor: Colors.byName('green-400')
      },
      data: [
          ['Jan', 35],
          ['Feb', 25],
          ['Mar', 45],
          ['Apr', 25],
          ['May', 30],
          ['Jun', 15]
      ]
  }];

  $('#flot-bar-chart').each(function() {
      var $el = $(this);
      if ($el.data('height')) $el.height($el.data('height'));
      // $el.plot(barChartData, barChartOptions);
      $el.plot(data, barChartOptions);

  });

  $("#flot-bar-chart").bind("plotclick", function (event, pos, item) {
      if (item) {
          getVoucher(item.series.label);
      }
  });

  $('#chart-category').html("");
  var i;
  for( i = 0; i < data.length; i++){
    var elemBody = "<em class='ion-record spr' style='color:"+data[i].bars.fillColor+"'></em>"
        + "<small class='va-middle'>"+data[i].label+"</small>";
    var elem = $("<span class='mr'></span>").html(elemBody);
    elem.appendTo('#chart-category');
  }
}

function initFlotLine(data){
  // line chart
  //
  //--------------------------------------------------------
      var lineData = data;
      var lineOptions = {
          series: {
              lines: {
                  show: true,
                  fill: 0.01
              },
              points: {
                  show: true,
                  radius: 4
              }
          },
          grid: {
              borderColor: 'rgba(162,162,162,.26)',
              borderWidth: 1,
              hoverable: true,
              backgroundColor: 'transparent'
          },
          tooltip: true,
          tooltipOpts: {
              content: function(label, x, y) {
                  return x + ' : ' + y;
              }
          },
          xaxis: {
              tickColor: 'rgba(162,162,162,.26)',
              font: {
                  color: Colors.byName('blueGrey-200')
              },
              mode: 'categories'
          },
          yaxis: {
              // position: (isRTL ? 'right' : 'left'),
              tickColor: 'rgba(162,162,162,.26)',
              font: {
                  color: Colors.byName('blueGrey-200')
              }
          },
          shadowSize: 0
      };

      $('#line-flotchart').plot(lineData, lineOptions);
      $("#flot-bar-chart").attr( "style", "display:none" );
      $("#line-flotchart").attr( "style", "display:block; padding: 0px; position: relative;" );
}
