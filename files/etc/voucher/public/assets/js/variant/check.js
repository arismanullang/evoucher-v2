$( document ).ready(function() {
  var id = findGetParameter('id');
  $("#variant-id").val(id);
  getVoucher(id);
  getPartner(id);

  $('#profileForm').submit(function(e) {
       e.preventDefault();
       e.returnValue = false;
  });
});

function getVoucher(id) {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/vouchers?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;
          if (arrData.length > 4){
            limit = 4;
            $("<div class='card-body pv0 text-right'><a href='/voucher/search?variant_id="+id+"' class='btn btn-flat btn-info'>View all</a></div>").appendTo('#cardVoucher');
          }

          for ( i = 0; i < limit; i++){
            var html = "<div class='mda-list-item-icon'><em class='ion-pricetag icon-2x'></em></div>"
            +  "<div class='mda-list-item-text'>"
            +  "<h3><a href='/voucher/check?id="+arrData[i].id+"'>"+arrData[i].voucher_code+"</a></h3>"
            +  "<div class='text-muted text-ellipsis'>Status "+arrData[i].state+"</div>"
            +"</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listVoucher');
          }

          getVariant(id, arrData.length);
        },
        error: function (data) {
          console.log(data.data);
          $("<div class='card-body text-center'>No Voucher Yet</div>").appendTo('#cardVoucher');
          getVariant(id, 0);
        }
    });
}

function getPartner(id) {
    console.log("Get Partner Data");

    var arrData = [];
    $.ajax({
        url: '/v1/api/get/partner?variant_id='+id+'&token='+token,
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
            +  "<div class='text-muted text-ellipsis'>"+arrData[i].serial_number+"</div>"
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

function getVariant(id, voucher) {
    console.log("Get Variant Data");
    console.log(voucher);
    var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    var arrData = [];
    $.ajax({
        url: '/v1/api/get/variant/'+id+'?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;

	  var date1 = result.start_date.substring(0, 10).split("-");
	  var date2 = result.end_date.substring(0, 10).split("-");
          var startDate = date1[2] + " " + months[parseInt(date1[1])-1] + " " + date1[0];
          var endDate = date2[2] + " " + months[parseInt(date2[1])-1] + " " + date2[0];

          var period = startDate + " - " + endDate;
	  var variantType = "Email Blast"
          var remainingVoucher = result.max_quantity_voucher;
	  var maxVoucher = result.max_quantity_voucher;
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
          $('#maxQuantityVoucher').html(maxVoucher);
          $('#voucherValue').html(result.discount_value);
          $('#period').html(period);
          $('#variantTnc').html(result.variant_tnc);
          $('#remainingVoucher').html(remainingVoucher);
        //   $('#variant-image').attr("src",result.image_url);
        }
    });
}

// $("#btnExport").click(function (e) {
//
// });

function ConvertToCSV(objArray) {
	var array = typeof objArray != 'object' ? JSON.parse(objArray) : objArray;
	var str = '';

	for (var i = 0; i < array.length; i++) {
		var line = '';
		for (var index in array[i]) {
			if (line != '') line += ','

			line += array[i][index];
		}

		str += line + '\n';
	}

	return str;
}

function generateVoucher() {
	var id = $('#variant-id').val();
	console.log("Get Variant Data");
	$.ajax({
		url: '/v1/voucher/generate/bulk?variant='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data);
			alert("Success");
		}
	});
}

function generateLink() {
	var id = $('#variant-id').val();
	$.ajax({
		url: '/v1/voucher/link?variant='+id+"&token="+token,
		type: 'get',
		success: function (data) {
			console.log(data);
			var csv = ConvertToCSV(data.data);
			alert(csv);
			console.log(csv);
			window.open(encodeURI('data:text/csv;charset=utf-8,' + csv));
		}
	});

}

function editVariant(){
  var id = findGetParameter("id");
  window.location = "/variant/update?id="+id;
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
