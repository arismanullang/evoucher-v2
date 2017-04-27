$(document).ready(function() {
  getVariant();
});

function getVariant() {
    console.log("Get Variant Data");
    $.ajax({
        url: '/v1/report/variant',
        type: 'get',
        success: function (data) {
          console.log(data.data.chart);
          var result = data.data.chart;
          initFlotBar("variant", result);

          console.log(data.data.data);
          arrData = data.data.data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            dataSet[i] = [
              arrData[i].month
              , arrData[i].total
              , arrData[i].username
            ];
          }
          console.log(dataSet);

          if ($.fn.DataTable.isDataTable("#variant-datatable")) {
            $('#variant-datatable').DataTable().clear().destroy();
          }

          var table = $('#variant-datatable').dataTable({
              data: dataSet,
              dom: 'Brt',
              buttons: [
                  'copy', 'csv', 'excel', 'pdf', 'print'
              ],
              columns: [
                  { title: "month" },
                  { title: "total" },
                  { title: "username" }
              ],
              oLanguage: {
                  sSearch: '<em class="ion-search"></em>',
                  sLengthMenu: '_MENU_ records per page',
                  info: 'Showing page _PAGE_ of _PAGES_',
                  zeroRecords: 'Nothing found - sorry',
                  infoEmpty: 'No records available',
                  infoFiltered: '(filtered from _MAX_ total records)',
                  oPaginate: {
                      sNext: '<em class="ion-ios-arrow-right"></em>',
                      sPrevious: '<em class="ion-ios-arrow-left"></em>'
                  }
              }
            });
        },
        error: function (data) {
          alert("Variant Not Found.");
        }
    });
}

function getVoucher(id) {
    console.log("Get Variant Data");
    $.ajax({
        url: '/v1/report/vouchers/variant?id='+id,
        type: 'get',
        success: function (data) {
          console.log(data.data.chart);
          var result = data.data.chart;
          // initFlotBar("voucher", result);
          initFlotStacked("voucher", result);

          console.log(data.data.data);
          arrData = data.data.data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            dataSet[i] = [
              arrData[i].id
              , arrData[i].variant_name
              , arrData[i].total
              , arrData[i].quota
              , arrData[i].created_by
              , arrData[i].state
            ];
          }
          console.log(dataSet);

          if ($.fn.DataTable.isDataTable("#voucher-datatable")) {
            $('#voucher-datatable').DataTable().clear().destroy();
          }

          var table = $('#voucher-datatable').dataTable({
              data: dataSet,
              dom: 'Brt',
              buttons: [
                  'copy', 'csv', 'excel', 'pdf', 'print'
              ],
              columns: [
                  { title: "id" },
                  { title: "variant_name" },
                  { title: "total" },
                  { title: "quota" },
                  { title: "created_by" }
              ],
              oLanguage: {
                  sSearch: '<em class="ion-search"></em>',
                  sLengthMenu: '_MENU_ records per page',
                  info: 'Showing page _PAGE_ of _PAGES_',
                  zeroRecords: 'Nothing found - sorry',
                  infoEmpty: 'No records available',
                  infoFiltered: '(filtered from _MAX_ total records)',
                  oPaginate: {
                      sNext: '<em class="ion-ios-arrow-right"></em>',
                      sPrevious: '<em class="ion-ios-arrow-left"></em>'
                  }
              }
            });

            $('#voucher').attr("style", "display:block");
        },
        error: function (data) {
          alert("Voucher Not Found.");
        }
    });
}

function initFlotBar(elem ,data){
  // Flot bar chart
  // ------------------
  var barChartOptions = {
      series: {
          bars: {
              show: true,
              fill: 1,
              fillColor: {
                  colors: [{
                      opacity: 1
                  }, {
                      opacity: 1
                  }]
              },
              barWidth: 0.3,
              lineWidth: 0,
              align: 'center'
          },
          highlightColor: 'rgba(255,255,255, 0.2)'
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
              return label + ' : ' + y;
          }
      },
      xaxis: {
          tickColor: 'transparent',
          mode: 'categories',
          font: {
              color: "black"
          }
      },
      yaxis: {
          tickColor: 'rgba(162,162,162,.26)',
          font: {
              color: "black"
          }
      },
      legend: {
          show: false
      },
      shadowSize: 0
  };

  var chartName = "#" + elem + "-flot-bar-chart";
  $(chartName).each(function() {
      var $el = $(this);
      if ($el.data('height')) $el.height($el.data('height'));
      // $el.plot(barChartData, barChartOptions);
      $el.plot(data, barChartOptions);

  });

  $(chartName).bind("plotclick", function (event, pos, item) {
    if(chartName == "#variant-flot-bar-chart"){
      if (item) {
          getVoucher(item.series.label);
      }
    }
  });

  var chartCategory = "#" + elem + "-chart-category";
  $(chartCategory).html("");
  var i;
  for( i = 0; i < data.length; i++){
    var elemBody = "<em class='ion-record spr' style='color:"+data[i].bars.fillColor+"'></em>"
        + "<small class='va-middle'>"+data[i].label+"</small>";
    var elem = $("<span class='mr'></span>").html(elemBody);
    elem.appendTo(chartCategory);
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

function initFlotStacked(elem, data){
  // Bar chart stacked
  // ------------------------

  var label = [];

  var data1 = data[0];
  var d1data = [];
  for (var i = 0; i < data1.data.length; i += 1) {
    if(!(data1.data[i][0] in label)){
      label.push(data1.data[i][0]);
    }
  	d1data.push([label.indexOf(data1.data[i][0]), data1.data[i][1]]);
  }

  var data1 = data[1];
  var d2data = [];
  for (var i = 0; i < data1.data.length; i += 1) {
    if(label.indexOf(data1.data[i][0]) == -1){
      label.push(data1.data[i][0]);
    }
  	d2data.push([label.indexOf(data1.data[i][0]), data1.data[i][1]]);
  }

 var stackedChartOptions = {
   series: {
        stack: true,
        bars: {
            align: 'center',
            lineWidth: 0,
            show: true,
            barWidth: 0.3,
            fill: true,
            fillColor: {
                colors: [{
                    opacity: 1
                }, {
                    opacity: 1
                }]
            }
        }
    },
    grid: {
        borderColor: 'rgba(162,162,162,.26)',
        borderWidth: 0,
        hoverable: true
    },
    xaxis: {
        tickColor: 'transparent',
        font: {
            color: "black" //Colors.byName('blueGrey-200')
        },
        mode: 'categories'
    },
    tooltip: true,
    tooltipOpts: {
        content: function getTooltip(label, x, y) {
            return label + ' : ' + y;
        }
    },
    yaxis: {
        // position: (isRTL ? 'right' : 'left'),
        tickColor: 'rgba(162,162,162,.26)',
        font: {
            color: "black"
        }
    },
    legend: {
        show: false
    },
    shadowSize: 0
 };

 var chartName = "#" + elem + "-flot-bar-chart";
 $(chartName).each(function() {
     var $el = $(this);
     if ($el.data('height'))
      $el.height($el.data('height'));

     $el.plot([ { label: "Distributed Voucher", data: d1data },{ label: "Remaining Voucher", data: d2data }], stackedChartOptions);
 });

 var chartCategory = "#" + elem + "-chart-category";
 $(chartCategory).html("");
 var i;
 for( i = 0; i < label.length; i++){
   var elemBody = "<span class='badge'>"+i+"</span>"
       + "<small class='va-middle'> &nbsp;"+label[i]+"</small>";
   var elem = $("<span class='mr'></span>").html(elemBody);
   elem.appendTo(chartCategory);
 }
}
