$( window ).ready(function() {
  var id = findGetParameter("id");
  searchById(id);
  getPartner(id);
  $("#image-url").change(function() {
	readURL(this);
  });


  // getPartner();
});

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		reader.onload = function (e) {
			$('#image-preview').attr('src', e.target.result);
		}

		reader.readAsDataURL(input.files[0]);
	}
}

function searchById(id) {

    var arrData = [];

    $.ajax({
        url: '/v1/ui/variant/detail?id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data);
          var variant = data.data[0];
	  $("#variant-id").val(id);
          $("#variant-name").val(variant.variant_name);
          $("#variant-type").val(variant.variant_type);
          $("#voucher-type").val(variant.voucher_type);
          $("#voucher-price").val(variant.voucher_price);
          $("#max-quantity-voucher").val(variant.max_quantity_voucher);
          $("#redemption-method").val(variant.redeemtion_method);
          $("#variant-valid-from").val(convertToDate(variant.start_date));
          $("#variant-valid-to").val(convertToDate(variant.end_date));
          $("#voucher-value").val(variant.discount_value);
          $("#variant-tnc").html(variant.variant_tnc);
          $("#variant-description").val(variant.variant_description);
          $("#start-hour").val(variant.start_hour);
          $("#end-hour").val(variant.end_hour);
	  $("#image-url-default").val(variant.image_url);
	  $("#voucher-valid-from").val(variant.valid_voucher_start);
	  $("#voucher-valid-to").val(variant.valid_voucher_end);

	  $("#voucher-price").attr("disabled","");
	  $("#max-quantity-voucher").attr("disabled","");
          $("#voucher-value").attr("disabled","");
          $("#start-hour").attr("disabled","");
          $("#end-hour").attr("disabled","");
          $("#variant-valid-from").attr("disabled","");
          $("#variant-valid-to").attr("disabled","");
	  $("#voucher-validity-type").attr("disabled","");
	  $("#voucher-valid-from").attr("disabled","");
	  $("#voucher-valid-to").attr("disabled","");

	  if(variant.voucher_lifetime != 0){
		$("#voucher-lifetime").attr("disabled","");
		$("#voucher-lifetime").val(variant.voucher_lifetime);
		$("#validity-lifetime").attr("style","display:block");
		$("#validity-date").attr("style","display:none");
		$("#voucher-valid-from").val("");
		$("#voucher-valid-to").val("");
	  	$("#voucher-validity-type").selectedIndex = 1;
	  	$("#voucher-validity-type").val("lifetime");
	  }
	  if($("#voucher-validity-type").val() == "lifetime"){
	    $("#validity-lifetime").attr("style","display:block");
	    $("#validity-date").attr("style","display:none");
	    $("#voucher-valid-from").val("");
	    $("#voucher-valid-to").val("");
	  } else if($("#voucher-validity-type").val() == "period"){
	    $("#validity-lifetime").attr("style","display:none");
	    $("#validity-date").attr("style","display:block");
	    $("#voucher-lifetime").val("");
	  } else {
	    $("#validity-lifetime").attr("style","display:none");
	    $("#validity-date").attr("style","display:none");
	    $("#voucher-valid-from").val("");
	    $("#voucher-valid-to").val("");
	    $("#voucher-lifetime").val("");
	  }

	  $("#redeem-validity-type").attr("disabled","");
	  if( $("#redeem-validity-type").val() == "all"){
	    $("#validity-day").attr("style","display:none");
	  } else if( $("#redeem-validity-type").value == "selected"){
	    $("#validity-day").attr("style","display:block");
	  } else {
	    $("#validity-day").attr("style","display:none");
	  }

	  $("#variant-type").attr("disabled","");
	  if($("#variant-type").val() == "bulk"){
	    $("#target").attr("style","display:block");
	    $("#max-quantity-voucher").attr("disabled","");
	    $("#max-usage-voucher").attr("disabled","");
	    $("#voucher-price").attr("disabled","");
	  } else{
	    $("#target").attr("style","display:none");
	    $("#max_quantity_voucher").removeAttr("disabled","");
	    $("#max_usage_voucher").removeAttr("disabled","");
	    $("#voucher_price").removeAttr("disabled","");
	  }


	  $('.summernote').each(function(){
		$(this).summernote({
			height: 380,
			placeholder: 'Any Message...'
		});
	  });
        }
    });
}

function send() {
  error = false;
  if($("#redeem-validity-type").val() == "all"){
    listDay = "all";
  } else if($("#redeem-validity-type").val() == "selected"){
    var li = $( "ul.select2-selection__rendered" ).find( "li" );

    if(li.length == 0 || parseInt($("#length").val()) < 8){
      error = true;
    }

    for (i = 0; i < li.length-1; i++) {
        var text = li[i].getAttribute("title");
        var value = $("option").filter(function() {
          return $(this).text() === text;
        }).first().attr("value");

        listDay = listDay + value+";";
    }
  }

  var listPartner = [];
  var li = $( "input[type=checkbox]:checked" );

  if(li.length == 0 || parseInt($("#length").val()) < 8){
	error = true;
  }

  for (i = 0; i < li.length; i++) {
	listPartner[i] = li[i].value;
  }

  var lifetime = 0;
  var periodStart = "";
  var periodEnd = "";

  if($("#voucher-validity-type").val() == "period"){
    lifetime = 0;
    periodStart = $("#voucher-valid-from").val();
    periodEnd = $("#voucher-valid-to").val();
  }else if($("#voucher-validity-type").val() == "lifetime"){
    lifetime = $("#voucher-lifetime").val();
    periodStart = "1001-01-01T00:00:00Z";
    periodEnd = "1001-01-01T00:00:00Z";
  }

    $('input[check="true"]').each(function() {
      if($(this).val() == ""){
        $(this).addClass("error");
        $(this).parent().closest('div').addClass("input-error");
        alert($(this).val());
        error = true;
      }
    });

    var str = $("#variant-tnc").summernote('code');
    var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');

    if(error){
      alert("Please check your input.");
      return
    }

    var formData = new FormData();
    var img = $('#image-url-default').val();
    if($('#image-url')[0].files[0] != null){

     formData.append('image-url', $('#image-url')[0].files[0]);

     jQuery.ajax({
         url:'/file/upload',
         type:"POST",
         processData: false,
         contentType: false,
         data: formData,
         success: function(data){
           console.log(data.data);
           img = data.data;

  	   var id = $("#variant-id").val();
	   var variant = {
		 variant_name: $("#variant-name").val(),
		 variant_type: $("#variant-type").find(":selected").val(),
		 voucher_type: $("#voucher-type").find(":selected").val(),
		 voucher_price: parseInt($("#voucher-price").val()),
		 max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
		 max_usage_voucher: 1,
		 allowAccumulative: $("#allow-accumulative").is(":checked"),
		 redeemtion_method: $("#redeemtion-method").find(":selected").val(),
		 start_date: $("#variant-valid-from").val(),
		 end_date: $("#variant-valid-to").val(),
		 start_hour: $("#start-hour").val(),
		 end_hour: $("#end-hour").val(),
		 discount_value: parseInt($("#voucher-value").val()),
		 image_url: img,
		 variant_tnc: tnc,
		 variant_description: $("#variant-description").val(),
		 validity_days: listDay,
		 valid_voucher_start: periodStart,
		 valid_voucher_end: periodEnd,
		 voucher_lifetime: parseInt(lifetime)
	   };

	   console.log(variant);

	   $.ajax({
		 url: '/v1/ui/variant/update?id='+id+'&type=detail&token='+token,
		 type: 'post',
		 dataType: 'json',
		 contentType: "application/json",
		 data: JSON.stringify(variant),
		 success: function () {
			 var partner = {
				 user: "user",
				 data: listPartner
			 };

			 $.ajax({
				 url: '/v1/ui/variant/update?id='+id+'&type=tenant&token='+token,
				 type: 'post',
				 dataType: 'json',
				 contentType: "application/json",
				 data: JSON.stringify(partner),
				 success: function () {
					 var id = findGetParameter("id");
					 window.location = "/variant/check?id="+id+"&token="+token;
				 }
			 });
		 }
	   });
         }
     });
    }else {
	    var id = $("#variant-id").val();
	    var variant = {
		    variant_name: $("#variant-name").val(),
		    variant_type: $("#variant-type").find(":selected").val(),
		    voucher_type: $("#voucher-type").find(":selected").val(),
		    voucher_price: parseInt($("#voucher-price").val()),
		    max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
		    max_usage_voucher: 1,
		    allowAccumulative: $("#allow-accumulative").is(":checked"),
		    redeemtion_method: $("#redeemtion-method").find(":selected").val(),
		    start_date: $("#variant-valid-from").val(),
		    end_date: $("#variant-valid-to").val(),
		    start_hour: $("#start-hour").val(),
		    end_hour: $("#end-hour").val(),
		    discount_value: parseInt($("#voucher-value").val()),
		    image_url: img,
		    variant_tnc: tnc,
		    variant_description: $("#variant-description").val(),
		    validity_days: listDay,
		    valid_voucher_start: periodStart,
		    valid_voucher_end: periodEnd,
		    voucher_lifetime: parseInt(lifetime)
	    };

	    console.log(variant);

	    $.ajax({
		    url: '/v1/ui/variant/update?id='+id+'&type=detail&token='+token,
		    type: 'post',
		    dataType: 'json',
		    contentType: "application/json",
		    data: JSON.stringify(variant),
		    success: function () {
			    var partner = {
				    user: "user",
				    data: listPartner
			    };

			    $.ajax({
				    url: '/v1/ui/variant/update?id='+id+'&type=tenant&token='+token,
				    type: 'post',
				    dataType: 'json',
				    contentType: "application/json",
				    data: JSON.stringify(partner),
				    success: function () {
					    var id = findGetParameter("id");
					    window.location = "/variant/check?id="+id+"&token="+token;
				    }
			    });
		    }
	    });
    }
}

function getPartner(id) {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/ui/partner/all',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;

        var i;
        for (i = 0; i < arrData.length; i++){
          var li = $("<div class='col-sm-4'></div>");
          var html = "<label class='checkbox-inline c-checkbox'>"
                    + "<input type='checkbox' value='"+arrData[i].id+"' text='"+arrData[i].partner_name+"'>"
                    + "<span class='ion-checkmark-round'></span>" + arrData[i].partner_name
                    + "</label>";
          li.html(html);
          li.appendTo('#partner-list');
        }

	$.ajax({
            url: '/v1/ui/partner/variant?variant_id='+id+'&token='+token,
            type: 'get',
            success: function (data) {
              var i;
   	      var y;
   	      var li = $( "input[type=checkbox]" );

      	      for (i = 0; i < li.length; i++) {

   	          var tempElem = li[i];
                  var arrData = data.data;
                  var limit = arrData.length;
                  for ( y = 0; y < limit; y++){
   		       if(tempElem.getAttribute("text") == arrData[y].partner_name){
   			       tempElem.checked = true;
   		       }
                  }
   	      }
            },
            error: function (data) {
              console.log(data.data);
              $("<div class='card-body text-center'>No Partner Found</div>").appendTo('#cardPartner');
            }
        });
      }
     });
}

function convertToDate(date){
  var string1 = date.split("T")[0];
  var string2 = string1.split("-");
  var result = string2[1] + "/" + string2[2] + "/" + string2[0];

  return result;
}

function convertToUpperCase(upperCase){
  var result = "";
  var firstChar = upperCase.charAt(0);
  upperCase = upperCase.replace(firstChar, firstChar.toUpperCase());
  result = upperCase;

  return result;
}

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        $('.select2').select2();
        $("#validity-day").attr("style","display:none");
        $("#collapseThree").removeClass("in");
        $("#collapseTwo").removeClass("in");
        $('.datepicker4').datepicker({
                container:'#example-datepicker-container-4',
                autoclose: true,
                startDate: 'd',
                setDate: new Date()
            });

        $('.datepicker3').datepicker({
                container:'#example-datepicker-container-3',
                autoclose: true,
                startDate: 'd',
                setDate: new Date()
            });
        $('#startDate').datepicker('update', new Date());
        $('#endDate').datepicker('update', '+1d');

        var cpInput = $('.clockpicker').clockpicker();
        // auto close picker on scroll
        $('main').scroll(function() {
            cpInput.clockpicker('hide');
        });


    }

})();
