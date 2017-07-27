$(document).ready(function() {
  getProgram();
});

function getProgram() {
    console.log("Get Program Data");
    $.ajax({
        url: '/v1/ui/program/all?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var programList = [];
          var result = data.data;
          var limit = 5;

          var totalVoucher = 0;
          var totalUsedVoucher = 0;
          var totalGeneratedVoucher = 0;
          var totalProgram = 0;
          for ( i = 0; i < result.length; i++){
            if( totalProgram < 5){

              var date = result[i].end_date.substring(0, 10).split("-");
              var dateEnd  = new Date(date[0], date[1]-1, date[2]);
              var dateEnd_ms = dateEnd.getTime();
              var dateNow_ms  = Date.now();
              var one_day = 1000*60*60*24;
              var diffNow = Math.round((dateEnd_ms-dateNow_ms)/one_day);

              var html = "<h5 class='mb-sm'><a href='/program/check?id="+result[i].id+"&token="+token+"'>"+result[i].name+"</a></h5>"
              + "<p class='text-muted'>End in "+diffNow+" days</p>";
              if(result[i].vouchers == null){
                html += "<p>No voucher generated</p>";
              } else {
              	var voucher = 0;
              	for(var y = 0; y < result[i].vouchers.length; y++){
				voucher += parseInt(result[i].vouchers[y].voucher);

			if(result[i].vouchers[y].state != "created"){
				totalUsedVoucher += parseInt(result[i].vouchers[y].voucher);
			}else{
				totalGeneratedVoucher += parseInt(result[i].vouchers[y].voucher);
			}
		}
                html += "<p>"+voucher+" vouchers have distributed. "+(result[i].max_quantity_voucher - voucher)+" vouchers left.</p>";
              }
              var li = $("<li class='list-group-item'></li>").html(html);
              if(dateEnd_ms > dateNow_ms){
                  li.appendTo('#upcomming-program');
                  totalProgram++;
              }
            }

            if(!(programList.includes(result[i].id))){
              programList.push(result[i].id);
              totalVoucher += parseInt(result[i].max_quantity_voucher);
            }

          }

          $("#total-program").html(result.length);
          $("#total-voucher").html(totalVoucher);
          $("#total-generated").html(totalGeneratedVoucher);
          $("#total-used").html(totalUsedVoucher);
        },
        error: function (data) {
          alert("Program Not Found.");
        }
    });
}

function addProgram(){
  window.location = "/program/create?token="+token;
}

(function() {
    'use strict';
    $(initDashboard);

    function initDashboard() {

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
                    return 'Visitors for ' + x + ' was ' + (y * 1000);
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
            $el.plot(barChartData, barChartOptions);
        });

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

        // Sparklines
        // -----------------

        var sparkValue1 = [4, 4, 3, 5, 3, 4, 6, 5, 3, 2, 3, 4, 6];
        var sparkValue2 = [2, 3, 4, 6, 5, 4, 3, 5, 4, 3, 4, 3, 4, 5];
        var sparkValue3 = [4, 4, 3, 5, 3, 4, 6, 5, 3, 2, 3, 4, 6];
        var sparkValue4 = [6, 5, 4, 3, 5, 4, 3, 4, 3, 4, 3, 2, 2];
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

        initSparkline($('#sparkline1'), sparkValue1, sparkOpts);
        initSparkline($('#sparkline2'), sparkValue2, sparkOpts);
        initSparkline($('#sparkline3'), sparkValue3, sparkOpts);
        initSparkline($('#sparkline4'), sparkValue4, sparkOpts);
        // call sparkline and mix options with data attributes
        function initSparkline(el, values, opts) {
            el.sparkline(values, $.extend(sparkOpts, el.data()));
        }

    }
})();
