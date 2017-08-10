$( document ).ready(function() {
  getTransactionByPartner("");
  getPartner();

  // $("#partner-id").change(function() {
	//   console.log($("#partner-id").value);
	//   console.log($("#partner-id").val());
	//   getTransactionByPartner($("#partner-id").val());
  // });
});

function getPartner() {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/ui/partner/all?token='+token,
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;

        var i;
        for (i = 0; i < arrData.length; i++){
          var li = $("<option value='"+arrData[i].name+"'>"+arrData[i].name+"</div>");
          li.appendTo('#partner-id');
        }
      }
  });
}

function getTransactionByPartner(partnerId) {
    console.log("Get Program Data");

    var arrData = [];
    $.ajax({
        url: '/v1/ui/transaction/partner?token='+token+'&partner='+partnerId,
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
	  var username = [];
	  var usernameExist = false;

          for ( i = 0; i < arrData.length; i++){
	    usernameExist = false;
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

	    if(username.length == 0){
		usernameExist = true;
	    }else{
	    	for( var y = 0; y < username.length; y++){
			if(username[y] == arrData[i].username.String.toUpperCase()){
				usernameExist = false;
				break;
			}
			usernameExist = true;
		}
	    }

	    if(usernameExist){
	    	username.push(arrData[i].username.String.toUpperCase());
	    }
	    var tempVoucherLength = arrData[i].voucher.length;
	    for ( y = 0; y < tempVoucherLength; y++){
		    var tempArray = [arrData[i].partner_name.toUpperCase()
			    , arrData[i].transaction_code
			    , arrData[i].program_name.toUpperCase()
			    , arrData[i].voucher[y].VoucherCode
			    //, addDecimalPoints(arrData[i].discount_value)
			    , date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0]
			    , date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0]
			    , cashoutDate
			    , cashoutCashier.toUpperCase()
			    , status.toUpperCase()
		    ];
		    dataSet.push(tempArray);
	    }
	    i += tempVoucherLength-1;
          }
          console.log(dataSet);
	  console.log(username);

	  for( y = 0; y < username.length; y++){
		  var li = $("<option value='"+username[y]+"'>"+username[y]+"</div>");
		  li.appendTo('#username');
	  }

          var table = $('#datatable1').dataTable({
              data: dataSet,
              dom: 'lBrtip',
              buttons: [
                  'copy', 'csv', 'excel', 'pdf', 'print'
              ],
              "order": [[ 8, "desc" ]],
              columns: [
                  { title: "PARTNER" },
                  { title: "TRANSACTION CODE" },
                  { title: "PROGRAM" },
                  { title: "VOUCHER" },
                  { title: "ISSUED" },
                  { title: "REDEEM" },
                  { title: "CASHOUT" },
                  { title: "USER" },
                  { title: "STATUS" }
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
	    for( i = 0; i < columnInputs.length; i++){
		if(columnInputs.get(i).tagName.toLowerCase() == "select"){
			columnInputs[i].onchange = function() {
				table.fnFilter(this.value, columnInputs.index(this));
			};
		}else{
			columnInputs[i].onkeyup = function() {
				table.fnFilter(this.value, columnInputs.index(this));
			};
		}
	    }
        }
    });
}
