$( document ).ready(function() {
  getVoucher();
});

function getVoucher() {
    console.log("Get Voucher Data");
    var holder = findGetParameter('holder');
    var arrData = [];
    $.ajax({
        url: '/v1/voucher?holder='+holder+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          arrData = data.data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            var date1 = arrData[i].valid_at.substring(0, 10).split("-");
            var date2 = arrData[i].expired_at.substring(0, 10).split("-");

            var dateStart  = new Date(date1[0], date1[1]-1, date1[2]);
            var dateEnd  = new Date(date2[0], date2[1]-1, date2[2]);
            var dateNow  = new Date("2017", "02", "02");

            var one_day = 1000*60*60*24;
            var dateStart_ms = dateStart.getTime();
            var dateEnd_ms = dateEnd.getTime();
            var dateNow_ms = dateNow.getTime();

            var diffNow = Math.round((dateEnd_ms-dateNow_ms)/one_day);
            var diffTotal = Math.round((dateEnd_ms-dateStart_ms)/one_day);
            var persen = diffNow / diffTotal * 100;
            console.log(dateStart + " " + dateEnd + " " + dateNow);
            console.log(diffNow + " " + diffTotal + " " + persen);

            diffNow = diffNow + " hari";

            if( persen < 0){
              diffNow = "Expired";
            }
            if( persen == Infinity){
              diffNow = "Expired";
            }
// variant jangan id, name aja
            dataSet[i] = [
              arrData[i].voucher_code
              , arrData[i].reference_no
              , arrData[i].variant_id
              , "<div class='progress'>"
                + "<div role='progressbar progress-bar-success' aria-valuenow='"+diffNow+"' aria-valuemin='0' aria-valuemax='"+diffTotal+"' style='width: "+persen+"%;' class='progress-bar'>"+diffNow+"</div>"
                + "</div>"
              , arrData[i].state
              ];
          }
          console.log(dataSet);

          if ($.fn.DataTable.isDataTable("#datatable1")) {
            $('#datatable1').DataTable().clear().destroy();
          }

          $('#datatable1').dataTable({
              data: dataSet,
              columns: [
                  { title: "Voucher Code" },
                  { title: "Reference No" },
                  { title: "Variant" },
                  { title: "Period" },
                  { title: "State" }
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
