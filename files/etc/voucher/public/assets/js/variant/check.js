$( document ).ready(function() {
  var id = findGetParameter('id');
  $("#variant-id").val(id);
  $("#token").val(token);
  getVoucher(id);
  getPartner(id);

  // $('#profileForm').submit(function(e) {
	//   e.preventDefault();
	//   e.returnValue = false;
  // });
});

function getVoucher(id) {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/ui/vouchers?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;
          var used = 0;
          var paid = 0;

          var dataSet = [];
          for ( i = 0; i < limit; i++){
	    if(arrData[i].state == "used") {
		used++;
	    }

	    if(arrData[i].state == "paid") {
	    	used++;
	    	paid++;
	    }

            var dateValid = new Date(arrData[i].valid_at);
            var dateExpired = new Date(arrData[i].expired_at);
	    dataSet[i] = [
		  arrData[i].voucher_code
		  , arrData[i].holder.toUpperCase()
		  , dateValid.toDateString().toUpperCase()
		  , dateExpired.toDateString().toUpperCase()
		  , arrData[i].state.toUpperCase()
	    ];
          }

	  if ($.fn.DataTable.isDataTable("#datatable1")) {
		$('#datatable1').DataTable().clear().destroy();
	  }

	  var table = $('#datatable1').dataTable({
		data: dataSet,
		dom: 'lBrtip',
		buttons: [
			'copy', 'csv', 'excel', 'pdf', 'print'
		],
		"order": [[ 4, "desc" ]],
		columns: [
			{ title: "VOUCHER" },
			{ title: "HOLDER" },
			{ title: "VALID" },
			{ title: "EXPIRED" },
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

          getVariant(id, arrData.length, used, paid);
        },
        error: function (data) {
          console.log(data.data);
          getVariant(id, 0, 0, 0);
        }
    });
}

function getPartner(id) {
    console.log("Get Partner Data");

    var arrData = [];
    $.ajax({
        url: '/v1/ui/partner/variant?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;

          for ( i = 0; i < limit; i++){
            var html = "<div class='mda-list-item-icon'><em class='ion-ios-person icon-2x'></em></div>"
            +  "<div class='mda-list-item-text'>"
            +  "<h3><a href='#'>"+arrData[i].partner_name+"</a></h3>"
            +  "<div class='text-muted text-ellipsis'>"+arrData[i].serial_number.String+"</div>"
            +"</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listPartner');
          }
        },
        error: function (data) {
          console.log(data.data);
          $("<div class='card-body text-center'>No Partner Found</div>").appendTo('#cardPartner');
        }
    });
}

function getVariant(id, voucher, used, paid) {
    console.log("Get Variant Data");
    console.log(voucher);
    var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    var arrData = [];
    $.ajax({
        url: '/v1/ui/variant/detail?id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data[0];

	  var date1 = result.start_date.substring(0, 10).split("-");
	  var date2 = result.end_date.substring(0, 10).split("-");
          var startDate = date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0];
          var endDate = date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0];

          var period = startDate + " - " + endDate;
	  var variantType = "Email Blast"
          var remainingVoucher = result.max_quantity_voucher;
          if( voucher != null){
			remainingVoucher = result.max_quantity_voucher - voucher;
	  }

	  if(result.variant_type != 'bulk'){
	        variantType = "Mobile App"
		$("#button-link").attr("style","display:none");
		$("#button-voucher").attr("style","display:none");
	  }

          $('#variantName').html(result.variant_name);
          $('#variantDescription').html(result.variant_description);
          $('#variantType').html(variantType);
          $('#voucherType').html(result.voucher_type);
          $('#conversionRate').html(result.voucher_price);
          $('#maxQuantityVoucher').html(result.max_quantity_voucher);
          $('#voucherValue').html(result.discount_value);
          $('#period').html(period);
          $('#variantTnc').html(result.variant_tnc);
          $('#remainingVoucher').html(remainingVoucher);
          $('#createdVoucher').html(voucher);
          $('#usedVoucher').html(used);
          $('#paidVoucher').html(paid);
        //   $('#variant-image').attr("src",result.image_url);
        }
    });
}

function generateVoucher() {
	var id = $('#variant-id').val();
	console.log("Get Variant Data");
	$.ajax({
		url: '/v1/ui/voucher/generate/bulk?variant='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data);
			location.reload();
		}
	});
}

function editVariant(){
  var id = findGetParameter("id");
  window.location = "/variant/update?id="+id+"&token="+token;
}
