$( window ).ready(function() {
  var id = findGetParameter("id");
  searchById(id);
  getPartner(id);
  $("#image-url").change(function() {
	readURL(this);
  });

  $("#all-tenant").change(function() {
	console.log("a");
	var lis = $( "input[class=partner]" );
	if($("#all-tenant").is(':checked')){
		for (var i = 0; i < lis.length; i++) {
			lis[i].checked = true;
		}
	} else{
		for (var i = 0; i < lis.length; i++) {
			lis[i].checked = false;
		}
	}
  });
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
        url: '/v1/ui/program/detail?id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data);
          var program = data.data[0];
	  $("#program-id").val(id);
          $("#program-name").val(program.name);
          $("#program-type").val(program.type);
          $("#voucher-type").val(program.voucher_type);
          $("#voucher-price").val(program.voucher_price);
          $("#max-quantity-voucher").val(program.max_quantity_voucher);
          $("#max-generate-voucher").val(program.max_generate_voucher);
          $("#max-redeem-voucher").val(program.max_redeem_voucher);
	  $("#redemption-method").val(program.redeem_method);
          $("#program-valid-from").val(convertToDate(program.start_date));
          $("#program-valid-to").val(convertToDate(program.end_date));
          $("#voucher-value").val(program.voucher_value);
          $("#program-tnc").html(program.tnc);
          $("#program-description").val(program.description);
          $("#start-hour").val(program.start_hour);
          $("#end-hour").val(program.end_hour);
	  $("#image-url-default").val(program.image_url);
	  $("#voucher-valid-from").val(program.valid_voucher_start);
	  $("#voucher-valid-to").val(program.valid_voucher_end);
	  $("#all-tenant").prop("checked", false);

	  $("#program-type").attr("disabled","");
	  $("#voucher-price").attr("disabled","");
	  $("#max-quantity-voucher").attr("disabled","");
          $("#voucher-value").attr("disabled","");
          $("#start-hour").attr("disabled","");
          $("#end-hour").attr("disabled","");
          $("#program-valid-from").attr("disabled","");
          $("#program-valid-to").attr("disabled","");
	  $("#voucher-validity-type").attr("disabled","");
	  $("#voucher-valid-from").attr("disabled","");
	  $("#voucher-valid-to").attr("disabled","");

	  if(program.voucher_lifetime != 0){
		$("#voucher-lifetime").attr("disabled","");
		$("#voucher-lifetime").val(program.voucher_lifetime);
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

	  $("#program-type").attr("disabled","");
	  if($("#program-type").val() == "bulk"){
	    $("#target").attr("style","display:block");
	    $("#conversion-row").attr("style","display:none");
	    $("#generate-row").attr("style","display:none");
	    $("#max-quantity-voucher").attr("disabled","");
	    $("#voucher-price").attr("disabled","");
	  } else{
	    $("#target").attr("style","display:none");
	    $("#conversion-row").attr("style","display:block");
	    $("#generate-row").attr("style","display:block");
	    $("#max_quantity_voucher").removeAttr("disabled","");
	    $("#voucher_price").removeAttr("disabled","");
	  }
	  if(program.allow_accumulative){
		  $("#allow-accumulative").attr("checked",true);
	  }

	  $(".select2").select2();
	  $('.summernote').each(function(){
		$(this).summernote({
			height: 380,
			placeholder: 'Any Message...',
			callbacks: {
				onPaste: function (e) {
					var bufferText = ((e.originalEvent || e).clipboardData || window.clipboardData).getData('Text');

					e.preventDefault();

					// Firefox fix
					setTimeout(function () {
						document.execCommand('insertText', false, bufferText);
					}, 10);
				}
			}
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
  var li = $( "input[class=partner]:checked" );

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

  var maxRedeem = parseInt($("#max-redeem-voucher").val());
  var maxGenerate = parseInt($("#max-generate-voucher").val());

  $('input[check="true"]').each(function() {
      if($("#program-type").val() == "bulk"){
	    if(this.getAttribute("id") == "max-quantity-voucher" || this.getAttribute("id") == "max-generate-voucher" || this.getAttribute("id") == "voucher-price" || this.getAttribute("id") == "max-redeem-voucher"){
		    maxGenerate = 1;
		    maxRedeem = 1;
		    return true;
	    }
      }

      if($(this).val() == ""){
        $(this).addClass("error");
        $(this).parent().closest('div').addClass("input-error");
        alert($(this).attr("id") + " : " + $(this).val());
        error = true;
      }
  });


  if (maxRedeem < 1 || maxGenerate < 1){
	error = true;
  }

  var str = $("#program-tnc").summernote('code');
  var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');

  if(!str.includes("<p>")){
    	tnc = '<p>'+tnc+'</p>';
  }

  if(error){
      alert("Please check your input.");
      return
  }

  var formData = new FormData();
  var img = $('#image-url-default').val();
  var redeem = $("#redemption-method").val();
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

  	   var id = $("#program-id").val();
	   var program = {
		 name: $("#program-name").val(),
		 type: $("#program-type").find(":selected").val(),
		 voucher_type: $("#voucher-type").find(":selected").val(),
		 voucher_price: parseInt($("#voucher-price").val()),
		 max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
		 max_redeem_voucher: maxRedeem,
		 max_generate_voucher: maxGenerate,
		 allow_accumulative: $("#allow-accumulative").is(":checked"),
		 redemption_method: redeem,
		 start_date: $("#program-valid-from").val(),
		 end_date: $("#program-valid-to").val(),
		 start_hour: $("#start-hour").val(),
		 end_hour: $("#end-hour").val(),
		 voucher_value: parseInt($("#voucher-value").val()),
		 image_url: img,
		 tnc: tnc,
		 description: $("#program-description").val(),
		 validity_days: listDay,
		 valid_voucher_start: periodStart,
		 valid_voucher_end: periodEnd,
		 voucher_lifetime: parseInt(lifetime)
	   };

	   console.log(program);

	   $.ajax({
		 url: '/v1/ui/program/update?id='+id+'&type=detail&token='+token,
		 type: 'post',
		 dataType: 'json',
		 contentType: "application/json",
		 data: JSON.stringify(program),
		 success: function () {
			 var partner = {
				 user: "user",
				 data: listPartner
			 };

			 $.ajax({
				 url: '/v1/ui/program/update?id='+id+'&type=tenant&token='+token,
				 type: 'post',
				 dataType: 'json',
				 contentType: "application/json",
				 data: JSON.stringify(partner),
				 success: function () {
					 var id = findGetParameter("id");
					 window.location = "/program/check?id="+id+"&token="+token;
				 }
			 });
		 }
	   });
         }
     });
  }else {
	    var id = $("#program-id").val();
	    var program = {
		    name: $("#program-name").val(),
		    type: $("#program-type").find(":selected").val(),
		    voucher_type: $("#voucher-type").find(":selected").val(),
		    voucher_price: parseInt($("#voucher-price").val()),
		    max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
		    max_redeem_voucher: maxRedeem,
		    max_generate_voucher: maxGenerate,
		    allow_accumulative: $("#allow-accumulative").is(":checked"),
		    redemption_method: redeem,
		    start_date: $("#program-valid-from").val(),
		    end_date: $("#program-valid-to").val(),
		    start_hour: $("#start-hour").val(),
		    end_hour: $("#end-hour").val(),
		    voucher_value: parseInt($("#voucher-value").val()),
		    image_url: img,
		    tnc: tnc,
		    description: $("#program-description").val(),
		    validity_days: listDay,
		    valid_voucher_start: periodStart,
		    valid_voucher_end: periodEnd,
		    voucher_lifetime: parseInt(lifetime)
	    };

	    console.log(program);

	    $.ajax({
		    url: '/v1/ui/program/update?id='+id+'&type=detail&token='+token,
		    type: 'post',
		    dataType: 'json',
		    contentType: "application/json",
		    data: JSON.stringify(program),
		    success: function () {
			    var partner = {
				    user: "user",
				    data: listPartner
			    };

			    $.ajax({
				    url: '/v1/ui/program/update?id='+id+'&type=tenant&token='+token,
				    type: 'post',
				    dataType: 'json',
				    contentType: "application/json",
				    data: JSON.stringify(partner),
				    success: function () {
					    var id = findGetParameter("id");
					    window.location = "/program/check?id="+id+"&token="+token;
				    }
			    });
		    }
	    });
  }
}

function getPartner(id) {
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
          var li = $("<div class='col-sm-4'></div>");
          var html = "<label class='checkbox-inline c-checkbox'>"
                    + "<input type='checkbox' class='partner' value='"+arrData[i].id+"' text='"+arrData[i].name+"'>"
                    + "<span class='ion-checkmark-round'></span>" + arrData[i].name
                    + "</label>";
          li.html(html);
          li.appendTo('#partner-list');
        }

	$.ajax({
            url: '/v1/ui/partner/program?program_id='+id+'&token='+token,
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
   		       if(tempElem.getAttribute("text") == arrData[y].name){
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
    	$("#collapseFour").removeClass("in");
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
