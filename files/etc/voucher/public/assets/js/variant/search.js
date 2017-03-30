$( document ).ready(function() {
  getVariant();
});

function getVariant() {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: '/v1/api/get/allVariant?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          arrData = data.data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            var date1 = arrData[i].StartDate.substring(0, 10).split("-");
            var date2 = arrData[i].EndDate.substring(0, 10).split("-");

            var dateStart  = new Date(date1[0], date1[1]-1, date1[2]);
            var dateEnd  = new Date(date2[0], date2[1]-1, date2[2]);
            var dateNow_ms  = Date.now();

            var one_day = 1000*60*60*24;
            var dateStart_ms = dateStart.getTime();
            var dateEnd_ms = dateEnd.getTime();
            // var dateNow_ms = dateNow.getTime();

            var diffNow = Math.round((dateEnd_ms-dateStart_ms)/one_day);
            var persen = 100;

            if(dateStart_ms < dateNow_ms){
              diffNow = Math.round((dateEnd_ms-dateNow_ms)/one_day);
              var diffTotal = Math.round((dateEnd_ms-dateStart_ms)/one_day);
              persen = diffNow / diffTotal * 100;
            }

            if(dateNow_ms > dateEnd_ms){
              persen = -1;
            }

            console.log(arrData[i].Id + " " + dateStart + " " + dateEnd);
            console.log(arrData[i].Id + " " + diffNow + " " + diffTotal + " " + persen);

            diffNow = diffNow + " hari";

            if( persen < 0){
              diffNow = "Expired";
            }

            dataSet[i] = [
              arrData[i].VariantName
              , arrData[i].VoucherPrice
              , "Rp " + addDecimalPoints(arrData[i].DiscountValue) + ",00"
              , (arrData[i].MaxVoucher - arrData[i].Voucher)
              , "<div class='progress'>"
                + "<div role='progressbar' aria-valuenow='"+diffNow+"' aria-valuemin='0' aria-valuemax='"+diffTotal+"' style='width: "+persen+"%;' class='progress-bar'>"+diffNow+"</div>"
                + "</div>"
              , "<button type='button' onclick='detail(\""+arrData[i].Id+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"+
              "<button type='button' onclick='edit(\""+arrData[i].Id+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em></button>"+
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
                  { title: "Program Name" },
                  { title: "Conversion Rate" },
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
        }
    });
}

function addDecimalPoints(value) {
    var input = " "+value;
    var inputValue = input.replace('.', '').split("").reverse().join(""); // reverse
    var newValue = '';
    for (var i = 0; i < inputValue.length; i++) {
        if (i % 3 == 0 && i != 0 && i != inputValue.length-1) {
            newValue += '.';
        }
        newValue += inputValue[i];
    }
    return newValue.split("").reverse().join("");
}

function edit(url){
  window.location = "/variant/update?id="+url;
}

function detail(url){
  window.location = "/variant/check?id="+url;
}

function addVariant(url){
  window.location = "/variant/create";
}

function deleteVariant(id) {
    console.log("Delete Variant");

    $.ajax({
        url: '/v1/delete/variant/'+id+'?token='+token,
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
