$( document ).ready(function() {
  getVariant();
});

function getVariant() {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: '/v1/ui/variant/all?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          arrData = data.data;
          var i;
          var dataSet = [];
          var dataId = [];
          var dataType = [];
          var dataStart = [];
          var dataEnd = [];
          var dataModified = [];
          var dataName = [];
          var dataPrice = [];
          var dataValue = [];
          var dataMax = [];
          var dataVoucher = [];
          var dataRedeem = [];
          var dataStatus = [];
	  var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
          for ( i = 0; i < arrData.length; i++){
            var tempIndex = dataId.indexOf(arrData[i].id);
            if( tempIndex == -1){

	      var date1 = arrData[i].start_date.substring(0, 10).split("-");
	      var date2 = arrData[i].end_date.substring(0, 10).split("-");
	      var date3 = arrData[i].created_at.substring(0, 10).split("-");
	      var date4 = arrData[i].updated_at.String.substring(0, 10).split("-");
              dataId.push(arrData[i].id);
              if(arrData[i].variant_type == "on-demand"){
		dataType.push("Mobile App");
	      }else{
              	dataType.push("Email Blast");
	      }
              dataStart.push(date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0]);
              dataEnd.push(date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0]);
              dataName.push(arrData[i].variant_name);
              dataPrice.push(arrData[i].voucher_price);
              dataValue.push(arrData[i].discount_value);
              dataMax.push(arrData[i].max_quantity_voucher);

              var created = 0;
	      var redeem = 0;

              if(arrData[i].vouchers != null){
		for(y = 0; y < arrData[i].vouchers.length; y++){
		  created += parseInt(arrData[i].vouchers[y].voucher);
		  if(arrData[i].vouchers[y].state != 'created'){
		  	redeem += parseInt(arrData[i].vouchers[y].voucher);
		        created += parseInt(arrData[i].vouchers[y].voucher);
		  }
		}
	      }

              dataVoucher.push(created);
              dataRedeem.push(redeem);

              if(arrData[i].status = 'created' ){
	              	var dateStart  = new Date(date1[0], date1[1]-1, date1[2]);
	              	var dateEnd  = new Date(date2[0], date2[1]-1, date2[2]);
			if(Date.now() < dateStart.getTime()){
				dataStatus.push("Not Active");
			}else if(Date.now() > dateStart.getTime() && Date.now() < dateEnd.getTime()){
				dataStatus.push("Active");
			}else if(Date.now() > dateEnd.getTime()){
				dataStatus.push("End");
			}
	      } else if(arrData[i].status = 'deleted'){
		      	dataStatus.push("Disabled");
	      }

	      if(arrData[i].updated_at.String != ""){
		      dataModified.push(date4[2] + " " + months[parseInt(date4[1])-1] + " " + date4[0]);
	      }else{
		      dataModified.push(date3[2] + " " + months[parseInt(date3[1])-1] + " " + date3[0]);
	      }
            }
            else{
              //dataVoucher[tempIndex] = parseInt(dataVoucher[tempIndex]) + parseInt(arrData[i].voucher);
            }
          }

          for ( i = 0; i < dataId.length; i++){
            var button = "<button type='button' onclick='detail(\""+dataId[i]+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"+
            "<button type='button' onclick='edit(\""+dataId[i]+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em></button>"+
            "<button type='button' value=\""+dataId[i]+"\" class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em></button>"

	    var avail = 0;
	    var rate = 0;
	    if(dataMax[i] != 0){
		avail = dataMax[i] - dataVoucher[i];
		rate = dataRedeem[i]/dataMax[i]*100;
	    }
            dataSet[i] = [
              dataName[i].toUpperCase()
	      , dataType[i].toUpperCase()
              , dataPrice[i] + " / " + addDecimalPoints(dataValue[i])
	      , dataStatus[i].toUpperCase()
              , dataStart[i].toUpperCase()
              , dataEnd[i].toUpperCase()
	      , dataModified[i].toUpperCase()
	      , dataMax[i]
              , avail
              , rate+"%"
              , button
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
              "order": [[ 6, "desc" ]],
              columns: [
                  { title: "PROGRAM" },
                  { title: "TYPE" },
                  { title: "CONVERSION </br> (POINT / CURRENCY)" },
	          { title: "STATUS" },
                  { title: "START" },
                  { title: "END" },
                  { title: "LAST MODIFIED" },
                  { title: "TOTAL </br> VOUCHER" },
                  { title: "AVAILABLE </br> VOUCHER" },
                  { title: "REDEEM RATE" },
                  { title: "ACTION"}
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

function edit(url){
  window.location = "/variant/update?id="+url+"&token="+token;
}

function detail(url){
  window.location = "/variant/check?id="+url+"&token="+token;
}

function addVariant(url){
  window.location = "/variant/create?token="+token;
}

function deleteVariant(id) {
    console.log("Delete Variant");

    $.ajax({
        url: '/v1/ui/variant/delete?id='+id+'&token='+token,
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
