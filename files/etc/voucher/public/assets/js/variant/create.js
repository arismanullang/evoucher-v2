$( window ).load(function() {
  getPartner();
});

function toTwoDigit(val){
  if (val < 10){
    return '0'+val;
  }
  else {
    return val;
  }
}

function send() {
  var token = findGetParameter("token");

  var i;
  var listPartner = [];
  var li = $( "ul.select2-selection__rendered" ).find( "li" );

  for (i = 0; i < li.length-1; i++) {
      var text = li[i].getAttribute("title");
      var value = $("option").filter(function() {
        return $(this).text() === text;
      }).first().attr("value");

      listPartner[i] = value;
  }

  var voucherFormat = {
    prefix: $("#prefix").val(),
    postfix: $("#postfix").val(),
    body: $("#body").val(),
    format_type: $("#voucherFormat").find(":selected").val(),
    length: parseInt($("#length").val())
  }

  var variant = {
      variant_name: $("#variantName").val(),
      variant_type: $("#variantType").find(":selected").val(),
      voucher_format: voucherFormat,
      voucher_type: $("#voucherType").find(":selected").val(),
      voucher_price: parseInt($("#voucherPrice").val()),
      max_quantity_voucher: parseInt($("#maxQuantityVoucher").val()),
      max_usage_voucher: parseInt($("#maxUsageVoucher").val()),
      allowAccumulative: $("#allowAccumulative").is(":checked"),
      redeem_method: $("#redeemtionMethod").find(":selected").val(),
      start_date: $("#startDate").val(),
      end_date: $("#endDate").val(),
      discount_value: parseInt($("#voucherValue").val()),
      image_url: $("#imageUrl").val(),
      variant_tnc: $("#variantTnc").val(),
      variant_description: $("#variantDescription").val(),
      validUsers: listPartner
    };

    console.log(variant);
    $.ajax({
       url: 'http://evoucher.elys.id:8889/create/variant?token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(variant),
       success: function () {
           alert("Variant created.");
       }
   });
}

function getPartner() {
    console.log("Get Partner Data");
    var token = findGetParameter("token");

    $.ajax({
      url: 'http://evoucher.elys.id:8889/get/partner?token='+token,
      type: 'get',
      success: function (data) {
        renderData(data);
      }
  });
}

function renderData(data) {
  console.log("Render Data");
  var arrData = [];
  arrData = data.data.Data;

  var i;
  for (i = 0; i < arrData.length; i++){
    var li = $("<option value='"+arrData[i].Id+"'>"+arrData[i].PartnerName+"</option>");
    li.appendTo('#variantPartners');
  }
}

function setDefaultValue() {
  $("#voucherFormat").find("option")[0].attr("selected","selected");
  $("#voucherType").find("option")[0].attr("selected","selected");
  $("#variantType").find("option")[0].attr("selected","selected");
  $("#redeemMethod").find("option")[0].attr("selected","selected");
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

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        $('.select2').select2();

        $('.datepicker4')
            .datepicker({
                container:'#example-datepicker-container-4'
            });
    }

})();
