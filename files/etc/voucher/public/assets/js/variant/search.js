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
          var dataId = [];
          var dataStart = [];
          var dataEnd = [];
          var dataModified = [];
          var dataName = [];
          var dataPrice = [];
          var dataValue = [];
          var dataMax = [];
          var dataVoucher = [];
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
              dataStart.push(date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0]);
              dataEnd.push(date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0]);
              dataName.push(arrData[i].variant_name);
              dataPrice.push(arrData[i].voucher_price);
              dataValue.push(arrData[i].discount_value);
              dataMax.push(arrData[i].max_quantity_voucher);
              dataVoucher.push(arrData[i].voucher);

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
              dataVoucher[tempIndex] = parseInt(dataVoucher[tempIndex]) + parseInt(arrData[i].voucher);
            }
          }

          for ( i = 0; i < dataId.length; i++){
        //     var date1 = dataStart[i].substring(0, 10).split("-");
        //     var date2 = dataEnd[i].substring(0, 10).split("-");
            //
        //     var dateStart  = new Date(date1[0], date1[1]-1, date1[2]);
        //     var dateEnd  = new Date(date2[0], date2[1]-1, date2[2]);
        //     var dateNow_ms  = Date.now();
            //
        //     var one_day = 1000*60*60*24;
        //     var dateStart_ms = dateStart.getTime();
        //     var dateEnd_ms = dateEnd.getTime();
        //     // var dateNow_ms = dateNow.getTime();
            //
        //     var diffNow = Math.round((dateEnd_ms-dateStart_ms)/one_day);
        //     var persen = 100;
            //
        //     if(dateStart_ms < dateNow_ms){
        //       diffNow = Math.round((dateEnd_ms-dateNow_ms)/one_day);
        //       var diffTotal = Math.round((dateEnd_ms-dateStart_ms)/one_day);
        //       persen = diffNow / diffTotal * 100;
        //     }
            //
        //     if(dateNow_ms > dateEnd_ms){
        //       persen = -1;
        //     }
            //
        //     console.log(dataId[i] + " " + dateStart + " " + dateEnd);
        //     console.log(dataId[i] + " " + diffNow + " " + diffTotal + " " + persen);
        //     var diffDay = diffNow;
        //     diffNow = diffNow + " hari";
            //
        //     if( persen < 0){
        //       diffNow = "Expired";
        //     }
            //
        //     if( diffDay == 30 && diffTotal > 30){
        //       diffNow = "Not start yet";
        //     }

            var button = "<button type='button' onclick='detail(\""+dataId[i]+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"+
            "<button type='button' onclick='edit(\""+dataId[i]+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em></button>"+
            "<button type='button' value=\""+dataId[i]+"\" class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em></button>"

            dataSet[i] = [
              dataName[i]
              , dataPrice[i] + " / " + addDecimalPoints(dataValue[i])
              , dataStart[i]
              , dataEnd[i]
	      , dataStatus[i]
	      , dataModified[i]
	      , dataMax[i]
              , (dataMax[i] - dataVoucher[i])
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
              "order": [[ 5, "desc" ]],
              columns: [
                  { title: "Program Name" },
                  { title: "Conversion </br> (point / currency)" },
                  { title: "Start Date" },
                  { title: "End Date" },
                  { title: "Status" },
                  { title: "Last Modified" },
                  { title: "Total Voucher" },
                  { title: "Available Voucher" },
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
            var inputSearchClass = 'datatable_input_col_search';
            var columnInputs = $('thead .' + inputSearchClass);

            columnInputs
                .keyup(function() {
                    table.fnFilter(this.value, columnInputs.index(this));
                });
        }
    });
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
