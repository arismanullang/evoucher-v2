$( document ).ready(function() {
  getPartner();

  $("#token").val(token);
  $("#voucher-validity-type").change(function() {
    if(this.value == "lifetime"){
      $("#validity-lifetime").attr("style","display:block");
      $("#validity-date").attr("style","display:none");
      $("#voucher-valid-from").val("");
      $("#voucher-valid-to").val("");
    } else if(this.value == "period"){
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
  });
  $("#redeem-validity-type").change(function() {
    if(this.value == "all"){
      $("#validity-day").attr("style","display:none");
    } else if(this.value == "selected"){
      $("#validity-day").attr("style","display:block");
    } else {
      $("#validity-day").attr("style","display:none");
    }
  });
  $("#variant-type").change(function() {
    if(this.value == "bulk"){
      $("#target").attr("style","display:block");
      $("#max-row").attr("style","display:none");
      $("#conversion-row").attr("style","display:none");
    } else{
      $("#target").attr("style","display:none");
      $("#max-row").attr("style","display:block");
      $("#conversion-row").attr("style","display:block");
    }
  });
  $("#image-url").change(function() {
    readURL(this);
  });

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

  if( $("#redeem-validity-type").val() == "all"){
   $("#validity-day").attr("style","display:none");
  } else if( $("#redeem-validity-type").val() == "selected"){
   $("#validity-day").attr("style","display:block");
  } else {
   $("#validity-day").attr("style","display:none");
  }

  if($("#variant-type").val() == "bulk"){
   $("#target").attr("style","display:block");
   $("#max-row").attr("style","display:none");
   $("#conversion-row").attr("style","display:none");
  } else{
   $("#target").attr("style","display:none");
   $("#max-row").attr("style","display:block");
   $("#conversion-row").attr("style","display:block");
  }
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

function addRule(){
  console.log("add");
  var body = "<td class='text-ellipsis td-index'>*</td>"
            + "<td class='text-ellipsis tnc'><div>"+$("#input-term-condition").val()+"</div></td>"
            + "<td><button type='button' onclick='removeElem(this)' class='btn btn-flat btn-sm btn-info pull-right'><em class='ion-close-circled'></em></button></td>";
  var li = $("<tr class='msg-display clickable'></tr>");
  li.html(body);
  li.appendTo('#list-rule');
  $('#input-term-condition').val('');
}

function toTwoDigit(val){
  if (val < 10){
    return '0'+val;
  }
  else {
    return val;
  }
}

function send() {
  error = false;
  var i;
  var listDay = "";
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
    periodStart = "01/01/0001";
    periodEnd = "01/01/0001";
  }

  var voucherFormat = {
    prefix: $("#prefix").val(),
    postfix: $("#postfix").val(),
    body: $("#body").val(),
    format_type: $("#voucher-format").find(":selected").val(),
    length: parseInt($("#length").val())
  }

  var tncTd = $('tr').find('td.tnc');
  var tnc = "";
  for (i = 0; i < tncTd.length; i++) {
    if(tncTd[i].innerHTML != ""){
      var decoded = $("<div/>").html((i+1) + ". " + tncTd[i].innerHTML).text();
      tnc += decoded + " <br>";
    }
  }

  var maxUsage = parseInt($("#max-usage-voucher").val());
  if($("#variant-type").val() == "bulk"){
	maxUsage = 1;
  }

  $('input[check="true"]').each(function() {
    if($("#variant-type").val() == "bulk"){
    	if(this.getAttribute("id") == "max-quantity-voucher" || this.getAttribute("id") == "max-usage-voucher" || this.getAttribute("id") == "voucher-price"){
		return true;
	}
    }
    if($(this).val() == ""){
      $(this).addClass("error");
      $(this).parent().closest('div').addClass("input-error");
      error = true;
    }

    if($(this).attr("id") == "length"){
      if(parseInt($(this).val()) < 8){
        error = true;
      }
    }
  });

  if(error){
    alert("Please check your input.");
    return
  }

  var formData = new FormData();
  var img = "https://storage.googleapis.com/e-voucher/Nd3QxH8El2Zuy12QhXs5Y305vPL4VZJJ.jpg";
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
         var variant = {
	       variant_name: $("#variant-name").val(),
	       variant_type: $("#variant-type").find(":selected").val(),
	       voucher_format: voucherFormat,
	       voucher_type: $("#voucher-type").find(":selected").val(),
	       voucher_price: parseInt($("#voucher-price").val()),
	       max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
	       max_usage_voucher: maxUsage,
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
	       validity_day: listDay,
	       valid_partners: listPartner,
	       valid_voucher_start: periodStart,
	       valid_voucher_end: periodEnd,
	       voucher_lifetime: parseInt(lifetime)
         };

         console.log(variant);

         $.ajax({
	       url: '/v1/create/variant?token='+token,
	       type: 'post',
	       dataType: 'json',
	       contentType: "application/json",
	       data: JSON.stringify(variant),
	       success: function (data) {
		       if($("#variant-type").find(":selected").val() == "bulk"){

			       var targets = new FormData();
			       targets.append('list-target', $("#list-target")[0].files[0]);

			       jQuery.ajax({
				       url:'/v1/upload/user?token='+token+'&variant-id='+data.data,
				       type:"POST",
				       processData: false,
				       contentType: false,
				       data: targets,
				       success: function(data){
					       console.log(data);
					       alert("Program created.");
					       window.location = "/variant/search";
				       }
			       });

		       }else{
			       alert("Program created.");
			       window.location = "/variant/search";
		       }
	       }
         });
       }
   });
  }else{
	  var variant = {
		  variant_name: $("#variant-name").val(),
		  variant_type: $("#variant-type").find(":selected").val(),
		  voucher_format: voucherFormat,
		  voucher_type: $("#voucher-type").find(":selected").val(),
		  voucher_price: parseInt($("#voucher-price").val()),
		  max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
		  max_usage_voucher: maxUsage,
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
		  validity_day: listDay,
		  valid_partners: listPartner,
		  valid_voucher_start: periodStart,
		  valid_voucher_end: periodEnd,
		  voucher_lifetime: parseInt(lifetime)
	  };
	  console.log(variant);

	  $.ajax({
		  url: '/v1/create/variant?token='+token,
		  type: 'post',
		  dataType: 'json',
		  contentType: "application/json",
		  data: JSON.stringify(variant),
		  success: function (data) {
			  if($("#variant-type").find(":selected").val() == "bulk"){

				  var targets = new FormData();
				  targets.append('list-target', $("#list-target")[0].files[0]);

				  jQuery.ajax({
					  url:'/v1/upload/user?token='+token+'&variant-id='+data.data,
					  type:"POST",
					  processData: false,
					  contentType: false,
					  data: targets,
					  success: function(data){
						  console.log(data);
						  alert("Program created.");
						  window.location = "/variant/search";
					  }
				  });

			  }else{
				  alert("Program created.");
				  window.location = "/variant/search";
			  }
		  }
	  });
  }
}

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
          var li = $("<div class='col-sm-4'></div>");
          var html = "<label class='checkbox-inline c-checkbox'>"
                    + "<input type='checkbox' value='"+arrData[i].id+"'>"
                    + "<span class='ion-checkmark-round'></span>" + arrData[i].partner_name
                    + "</label>";
          li.html(html);
          li.appendTo('#partner-list');
        }
      }
  });
}

function removeElem(elem){
  console.log("remove");
  $(elem).parent().closest('tr').remove();
}

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        $('.select2').select2();
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
