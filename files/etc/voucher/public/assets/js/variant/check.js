$( document ).ready(function() {
  getVariant();
});

var user = localStorage.getItem("user");
var token = localStorage.getItem(user);

function getVariant() {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: 'http://voucher.apps.id:8889/v1/api/get/allVariant?token='+token+'&user='+user,
        type: 'get',
        success: function (data) {
          console.log(data.data.Data);
          arrData = data.data.Data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            var date1 = arrData[i].StartDate.substring(0, 10).split("-");
            var date2 = arrData[i].EndDate.substring(0, 10).split("-");

            var dateStart  = new Date(date1[0], date1[1]-1, date1[2]);
            var dateEnd  = new Date(date2[0], date2[1]-1, date2[2]);

            var one_day = 1000*60*60*24;
            var dateStart_ms = dateStart.getTime();
            var dateEnd_ms = dateEnd.getTime();

            var diff = Math.round((dateEnd_ms-dateStart_ms)/one_day);
            var persen = diff*100;
            console.log(dateStart + " " + dateEnd);
            console.log(diff)

            dataSet[i] = [
              arrData[i].VariantName
              , arrData[i].VoucherPrice
              , arrData[i].DiscountValue
              , (arrData[i].MaxVoucher - arrData[i].Voucher)
              , "<div class='progress'>"
                + "<div role='progressbar' aria-valuenow='"+diff+"' aria-valuemin='0' aria-valuemax='"+diff+"' style='width: "+persen+"%;' class='progress-bar'>"+diff+" hari</div>"
                + "</div>"
              , "<button type='button' onclick='goTo(\""+arrData[i].Id+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em></button>"+
              "<button type='button' value=\""+arrData[i].Id+"\" class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em></button>"
            ];
          }
          console.log(dataSet);

          if ($.fn.DataTable.isDataTable("#datatable1")) {
            $('#datatable1').DataTable().clear().destroy();
          }

          $('#datatable1').dataTable({
              data: dataSet,
              columns: [
                  { title: "Variant Name" },
                  { title: "Voucher Price" },
                  { title: "Voucher Value" },
                  { title: "Remaining Voucher" },
                  { title: "Period" },
                  { title: "Action"}
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

          $('.ui-slider-values').each(function() {
              var slider = this;

              noUiSlider.create(slider, {
                  start: [0, 40],
                  connect: true,
                  // direction: 'rtl',
                  behaviour: 'tap-drag',
                  range: {
                      'min': 0,
                      'max': 100
                  }
              });

              slider.noUiSlider.on('slide', updateValues);
              updateValues();

              function updateValues() {
                  var values = slider.noUiSlider.get();
                  // Connecto to live values
                  $('.ui-slider-value-upper').html(values[1]);
              }
          });
        }
    });
}

function renderData(data) {
  console.log("Render Data");
  var arrData = [];
  arrData = data.data.Data;

  var i;
  for (i = 0; i < arrData.length; i++){
    var tr = $("<tr class='gradeA' role='row'></tr>");
    var td = "<td>"+arrData[i].Id+"</td>";
    td += "<td>"+arrData[i].VariantName+"</td>";
    td += "<td>"+arrData[i].VoucherPrice+"</td>";
    td += "<td>"+arrData[i].DiscountValue+"</td>";
    td += "<td>"+arrData[i].MaxVoucher+"</td>";
    tr.html(td);
    tr.appendTo('tbody');
  }
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
    .substr(1)
        .split("&")
        .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    });
    return result;
}

function goTo(url){
  window.location = "http://voucher.apps.id:8889/variant/update?id="+url;
}

function deleteVariant(id) {
    console.log("Delete Variant");

    $.ajax({
        url: 'http://voucher.apps.id:8889/delete/variant/'+id+'?token='+token+'&user='+user,
        type: 'get',
        success: function (data) {
          getVariant();
        }
    });
}



(function() {
    'use strict';

    $(runSweetAlert);
    //onclick='deleteVariant(\""+arrData[i].Id+"\")'
    function runSweetAlert() {
        $(document).on('click', '.swal-demo4', function(e) {
            e.preventDefault();
            console.log(e.target.value);
            swal({
                    title: 'Are you sure?',
                    text: 'Do you want delete variant?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes, delete it!',
                    closeOnConfirm: false
                },
                function() {
                    swal('Deleted!', 'Delete success.', deleteVariant(e.target.value));
                });

        });
    }

})();

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        // UI SLider (noUiSlider)
        $('.ui-slider').each(function() {

            noUiSlider.create(this, {
                start: $(this).data('start'),
                connect: 'lower',
                range: {
                    'min': 0,
                    'max': 100,
                }
            });
        });

        // Range selectable
        $('.ui-slider-range').each(function() {
            noUiSlider.create(this, {
                start: [25, 75],
                range: {
                    'min': 0,
                    'max': 100
                },
                connect: true
            });

        });

        // Live Values
        $('.ui-slider-values').each(function() {
            var slider = this;

            noUiSlider.create(slider, {
                start: [0, 40],
                connect: true,
                // direction: 'rtl',
                behaviour: 'tap-drag',
                range: {
                    'min': 0,
                    'max': 100
                }
            });

            slider.noUiSlider.on('slide', updateValues);
            updateValues();

            function updateValues() {
                var values = slider.noUiSlider.get();
                // Connecto to live values
                $('#ui-slider-value-lower').html(values[0]);
                $('#ui-slider-value-upper').html(values[1]);
            }
        });
    }

})();
