$(document).ready(function() {
  getVariant();

  var sparkValue1 = [4, 4, 3, 5, 3, 4, 6, 5, 3, 2, 3, 4, 6];
  var sparkValue2 = [2, 3, 4, 6, 5, 4, 3, 5, 4, 3, 4, 3, 4, 5];
  var sparkValue3 = [4, 4, 3, 5, 3, 4, 6, 5, 3, 2, 3, 4, 6];
  var sparkValue4 = [6, 5, 4, 3, 5, 4, 3, 4, 3, 4, 3, 2, 2];

  initSparkline($('#sparkline1'), sparkValue1);
  initSparkline($('#sparkline2'), sparkValue2);
  initSparkline($('#sparkline3'), sparkValue3);
  initSparkline($('#sparkline4'), sparkValue4);

  initFlotSplineData();
  initBar();
  initFlotBar();
  // $.get('/assets/server/chart/line.json', function(data) {
  //   initFlotLine(data);
  // });
});

function getVariant() {
    console.log("Get Variant Data");
    $.ajax({
        url: '/v1/report?id=DwXt2g5c',
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;

          console.log(result);
          initFlotLine(result);
        },
        error: function (data) {
          alert("Variant Not Found.");
        }
    });
}

function bubbleSorts(items, items2) {
    var length = items.length;
    for (var i = (length - 1); i >= 0; i--) {
        //Number of passes
        for (var j = (length - i); j > 0; j--) {
            //Compare the adjacent positions
            if (items[j] > items[j - 1]) {
                //Swap the numbers
                var tmp = items[j];
                items[j] = items[j - 1];
                items[j - 1] = tmp;

                var tmp = items2[j];
                items2[j] = items2[j - 1];
                items2[j - 1] = tmp;
            }
        }
    }
}

function bubbleSort(a, b){
    var swapped;
    do {
        swapped = false;
        for (var i=0; i < a.length-1; i++) {
            if (a[i] > a[i+1]) {
                var temp = a[i];
                a[i] = a[i+1];
                a[i+1] = temp;

                var tempB = b[i];
                b[i] = b[i+1];
                b[i+1] = temp;
                swapped = true;
            }
        }
    } while (swapped);
}

function addVariant(){
  window.location = "/variant/create";
}

function initFlotSplineData(){
  // Main Flot chart
  var splineData = [{
      'label': 'Clicks',
      'color': Colors.byName('purple-300'),
      data: [
          ['1', 40],
          ['2', 50],
          ['3', 40],
          ['4', 50],
          ['5', 66],
          ['6', 66],
          ['7', 76],
          ['8', 96],
          ['9', 90],
          ['10', 105],
          ['11', 125],
          ['12', 135]

      ]
  }, {
      'label': 'Unique',
      'color': Colors.byName('green-400'),
      data: [
          ['1', 30],
          ['2', 40],
          ['3', 20],
          ['4', 40],
          ['5', 80],
          ['6', 90],
          ['7', 70],
          ['8', 60],
          ['9', 90],
          ['10', 150],
          ['11', 130],
          ['12', 160]
      ]
  }, {
      'label': 'Recurrent',
      'color': Colors.byName('blue-500'),
      data: [
          ['1', 10],
          ['2', 20],
          ['3', 10],
          ['4', 20],
          ['5', 6],
          ['6', 10],
          ['7', 32],
          ['8', 26],
          ['9', 20],
          ['10', 35],
          ['11', 30],
          ['12', 56]

      ]
  }];
  var splineOptions = {
      series: {
          lines: {
              show: false
          },
          points: {
              show: false,
              radius: 3
          },
          splines: {
              show: true,
              tension: 0.39,
              lineWidth: 5,
              fill: 1,
              fillColor: Colors.byName('primary')
          }
      },
      grid: {
          borderColor: '#eee',
          borderWidth: 0,
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
          tickColor: 'transparent',
          mode: 'categories',
          font: {
              color: Colors.byName('blueGrey-200')
          }
      },
      yaxis: {
          show: false,
          min: 0,
          max: 220, // optional: use it for a clear representation
          tickColor: 'transparent',
          font: {
              color: Colors.byName('blueGrey-200')
          },
          //position: 'right' or 'left',
          tickFormatter: function(v) {
              return v /* + ' visitors'*/ ;
          }
      },
      shadowSize: 0
  };

  $('#flot-main-spline').each(function() {
      var $el = $(this);
      if ($el.data('height')) $el.height($el.data('height'));
      $el.plot(splineData, splineOptions);
  });
}

function initBar(){
  // Bar chart stacked
  // ------------------------
  var stackedChartData = [{
      data: [
          [1, 45],
          [2, 42],
          [3, 45],
          [4, 43],
          [5, 45],
          [6, 47],
          [7, 45],
          [8, 42],
          [9, 45],
          [10, 43]
      ]
  }, {
      data: [
          [1, 35],
          [2, 35],
          [3, 17],
          [4, 29],
          [5, 10],
          [6, 7],
          [7, 35],
          [8, 35],
          [9, 17],
          [10, 29]
      ]
  }];

  var stackedChartOptions = {
      bars: {
          show: true,
          fill: true,
          barWidth: 0.3,
          lineWidth: 1,
          align: 'center',
          // order : 1,
          fillColor: {
              colors: [{
                  opacity: 1
              }, {
                  opacity: 1
              }]
          }
      },
      colors: [Colors.byName('blue-100'), Colors.byName('blue-500')],
      series: {
          shadowSize: 3
      },
      xaxis: {
          show: true,
          position: 'bottom',
          ticks: 10,
          font: {
              color: Colors.byName('blueGrey-200')
          }
      },
      yaxis: {
          show: false,
          min: 0,
          max: 60,
          font: {
              color: Colors.byName('blueGrey-200')
          }
      },
      grid: {
          hoverable: true,
          clickable: true,
          borderWidth: 0,
          color: 'rgba(120,120,120,0.5)'
      },
      tooltip: true,
      tooltipOpts: {
          content: 'Value %x.0 is %y.0',
          defaultTheme: false,
          shifts: {
              x: 0,
              y: -20
          }
      }
  };

  $('#flot-stacked-chart').each(function() {
      var $el = $(this);
      if ($el.data('height')) $el.height($el.data('height'));
      $el.plot(stackedChartData, stackedChartOptions);
  });
}

function initFlotBar(){
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
      $el.plot(barChartData, barChartOptions);

  });

  var i;
  for( i = 0; i < barChartData.length; i++){
    var elemBody = "<em class='ion-record spr' style='color:"+barChartData[i].bars.fillColor+"'></em>"
        + "<small class='va-middle'>"+barChartData[i].label+"</small>";
    var elem = $("<span class='mr'></span>").html(elemBody);
    elem.appendTo('#chart-category');
  }
}

function initSmallFlot(){
  // Small flot chart
  // ---------------------
  var chartTaskData = [{
      'label': 'Total',
      color: Colors.byName('primary'),
      data: [
          ['Jan', 14],
          ['Feb', 14],
          ['Mar', 12],
          ['Apr', 16],
          ['May', 13],
          ['Jun', 14],
          ['Jul', 19]
          //4, 4, 3, 5, 3, 4, 6
      ]
  }];
  var chartTaskOptions = {
      series: {
          lines: {
              show: false
          },
          points: {
              show: false
          },
          splines: {
              show: true,
              tension: 0.4,
              lineWidth: 3,
              fill: 1
          },
      },
      legend: {
          show: false
      },
      grid: {
          show: false,
      },
      tooltip: true,
      tooltipOpts: {
          content: function(label, x, y) {
              return x + ' : ' + y;
          }
      },
      xaxis: {
          tickColor: '#fcfcfc',
          mode: 'categories'
      },
      yaxis: {
          min: 0,
          max: 30, // optional: use it for a clear representation
          tickColor: '#eee',
          //position: 'right' or 'left',
          tickFormatter: function(v) {
              return v /* + ' visitors'*/ ;
          }
      },
      shadowSize: 0
  };

  $('#flot-task-chart').each(function() {
      var $el = $(this);
      if ($el.data('height')) $el.height($el.data('height'));
      $el.plot(chartTaskData, chartTaskOptions);
  });
}

function initSparkline(el, values) {
  // Sparklines
  // -----------------
  var sparkOpts = {
      type: 'line',
      height: 20,
      width: '70',
      lineWidth: 2,
      valueSpots: {
          '0:': Colors.byName('blue-700'),
      },
      lineColor: Colors.byName('blue-700'),
      spotColor: Colors.byName('blue-700'),
      fillColor: 'transparent',
      highlightLineColor: Colors.byName('blue-700'),
      spotRadius: 0
  };
  el.sparkline(values, $.extend(sparkOpts, el.data()));
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
}
