$( window ).ready(function() {
  getVariant();
});

var user = localStorage.getItem("user");
var token = localStorage.getItem(user);

function getVariant() {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: 'http://evoucher.elys.id:8889/get/allVariant?token='+token+'&user='+user,
        type: 'get',
        success: function (data) {
          console.log(data.data.Data);
          arrData = data.data.Data;
          var i;
          var dataSet = [];
          for ( i = 0; i < arrData.length; i++){
            dataSet[i] = [
              arrData[i].Id
              , arrData[i].VariantName
              , arrData[i].VoucherPrice
              , arrData[i].DiscountValue
              , arrData[i].MaxVoucher
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
                  { title: "Id" },
                  { title: "Variant Name" },
                  { title: "Voucher Price" },
                  { title: "Discount Value" },
                  { title: "Max Voucher" },
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
  window.location = "http://evoucher.elys.id:8889/variant/update?id="+url;
}

function deleteVariant(id) {
    console.log("Delete Variant");

    $.ajax({
        url: 'http://evoucher.elys.id:8889/delete/variant/'+id+'?token='+token+'&user='+user,
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
