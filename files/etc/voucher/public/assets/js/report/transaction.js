$( document ).ready(function() {
  getTransaction();
  getPartner();

  $("#partner-id").change(function() {
	  console.log($("#partner-id").value);
	  console.log($("#partner-id").val());
	  getTransactionByPartner($("#partner-id").val());
  });
});

function getPartner() {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/get/partner',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;

        var i;
        for (i = 0; i < arrData.length; i++){
          var li = $("<option value='"+arrData[i].id+"'>"+arrData[i].partner_name+"</div>");
          li.appendTo('#partner-id');
        }
      }
  });
}

function getTransaction() {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: '/v1/get/transaction?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
	  var arrData = data.data;
          var i;
	  var dataSet = [];
	  var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

          for ( i = 0; i < arrData.length; i++){
	    var date1 = arrData[i].issued.substring(0, 10).split("-");
	    var date2 = arrData[i].redeemed.substring(0, 10).split("-");
	    var date3 = arrData[i].cashout.String.substring(0, 10).split("-");
	    var cashoutDate = date3[2] + " " + months[parseInt(date3[1])-1] + " " + date3[0];
	    var cashoutCashier = arrData[i].username.String;
	    var status = "Paid"
	    if( arrData[i].state == "used"){
		    cashoutDate = "-";
		    cashoutCashier = "-";
		    status = "Pending";
	    }
            dataSet[i] = [
              arrData[i].partner_name
              , arrData[i].voucher
              , addDecimalPoints(arrData[i].discount_value)
              , date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0]
	      , date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0]
	      , cashoutDate
              , cashoutCashier
	      , status
            ];
          }
          console.log(dataSet);

          if ($.fn.DataTable.isDataTable("#datatable1")) {
            $('#datatable1').DataTable().clear().destroy();
          }

          var table = $('#datatable1').dataTable({
              data: dataSet,
              dom: 'lBrtip',
              buttons: [
                  'copy', 'csv', 'excel', 'pdf', 'print'
              ],
              "order": [[ 0, "asc" ]],
              columns: [
                  { title: "Partner Name" },
                  { title: "Voucher Code" },
                  { title: "Voucher Value" },
                  { title: "Issued Date" },
                  { title: "Redeem Date" },
                  { title: "Cashout Date" },
                  { title: "Cashier" },
                  { title: "Status" }
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
            var inputSearchClass = 'datatable_input_col_search';
            var columnInputs = $('thead .' + inputSearchClass);

            columnInputs.keyup(function() {
                    table.fnFilter(this.value, columnInputs.index(this));
                });

        }
    });
}

function getTransactionByPartner(partnerId) {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: '/v1/get/transaction/partner?token='+token+'&partner='+partnerId,
        type: 'get',
        success: function (data) {
          if ($.fn.DataTable.isDataTable("#datatable1")) {
            $('#datatable1').DataTable().clear().destroy();
          }
          console.log(data.data);
	  var arrData = data.data;
          var i;
	  var dataSet = [];
	  var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

          for ( i = 0; i < arrData.length; i++){
	    var date1 = arrData[i].issued.substring(0, 10).split("-");
	    var date2 = arrData[i].redeemed.substring(0, 10).split("-");
	    var date3 = arrData[i].cashout.String.substring(0, 10).split("-");
	    var cashoutDate = date3[2] + " " + months[parseInt(date3[1])-1] + " " + date3[0];
	    var cashoutCashier = arrData[i].username.String;
	    var status = "Paid"
	    if( arrData[i].state == "used"){
		    cashoutDate = "-";
		    cashoutCashier = "-";
		    status = "Pending";
	    }
            dataSet[i] = [
              arrData[i].partner_name
              , arrData[i].voucher
              , addDecimalPoints(arrData[i].discount_value)
              , date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0]
	      , date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0]
	      , cashoutDate
              , cashoutCashier
	      , status
            ];
          }
          console.log(dataSet);

          var table = $('#datatable1').dataTable({
              data: dataSet,
              dom: 'lBrtip',
              buttons: [
                  'copy', 'csv', 'excel', 'pdf', 'print'
              ],
              "order": [[ 0, "asc" ]],
              columns: [
                  { title: "Partner Name" },
                  { title: "Voucher Code" },
                  { title: "Voucher Value" },
                  { title: "Issued Date" },
                  { title: "Redeem Date" },
                  { title: "Cashout Date" },
                  { title: "Cashier" },
                  { title: "Status" }
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
            var inputSearchClass = 'datatable_input_col_search';
            var columnInputs = $('thead .' + inputSearchClass);

            columnInputs.keyup(function() {
                    table.fnFilter(this.value, columnInputs.index(this));
                });

        }
    });
}
