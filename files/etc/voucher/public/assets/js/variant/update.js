var id = findGetParameter("id");
var user = localStorage.getItem("user");
var token = localStorage.getItem(user);

$( window ).ready(function() {
  searchById(id);
  // getPartner();
});

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

function searchById(id) {

    var arrData = [];

    $.ajax({
        url: '/v1/api/get/variant/'+id+'?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data);
          var variant = data.data;

          $("#variantName").val(variant.VariantName);
          $("#variantType").val(variant.VariantType);
          $("#select2-variantType-container").html(convcertToUpperCase(variant.VariantType));
          $("#voucherType").val(variant.VoucherType);
          $("#select2-voucherType-container").html(convcertToUpperCase(variant.VoucherType));
          $("#voucherPrice").val(variant.VoucherPrice);
          $("#maxQuantityVoucher").val(variant.MaxQuantityVoucher);
          $("#maxUsageVoucher").val(variant.MaxUsageVoucher);
          $("#redeemtionMethod").val(variant.RedeemtionMethod);
          $("#select2-redeemtionMethod-container").html(convcertToUpperCase(variant.RedeemtionMethod));
          $("#startDate").val(convertToDate(variant.StartDate));
          $("#endDate").val(convertToDate(variant.EndDate));
          $("#voucherValue").val(variant.DiscountValue);
          $("#imageUrl").val(variant.ImgUrl);
          $("#variantTnc").val(variant.VariantTnc);
          $("#variantDescription").val(variant.VariantDescription);

        }
    });
}

function send() {
    var variant = {
      variant_name: $("#variantName").val(),
      variant_type: $("#variantType").val(),
      voucher_type: $("#voucherType").val(),
      voucher_price: parseInt($("#voucherPrice").val()),
      max_quantity_voucher: parseInt($("#maxQuantityVoucher").val()),
      max_usage_voucher: parseInt($("#maxUsageVoucher").val()),
      redeemtion_method: $("#redeemtionMethod").val(),
      start_date: $("#startDate").val(),
      end_date: $("#endDate").val(),
      discount_value: parseInt($("#voucherValue").val()),
      image_url: $("#imageUrl").val(),
      variant_tnc: $("#variantTnc").val(),
      variant_description: $("#variantDescription").val()
    };

    console.log(variant);

    $.ajax({
        url: '/v1/update/variant/'+id+'?token='+token,
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(variant),
        success: function () {
            window.location = "/variant/search";
        }
    });
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
          var li = $("<option value='"+arrData[i].Id+"'>"+arrData[i].PartnerName+"</option>");
          li.appendTo('#variantPartners');
        }
      }
  });
}

function convertToDate(date){
  var string1 = date.split("T")[0];
  var string2 = string1.split("-");
  var result = string2[1] + "/" + string2[2] + "/" + string2[0];

  return result;
}

function convcertToUpperCase(upperCase){
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

        $('.datepicker4')
            .datepicker({
                container:'#example-datepicker-container-4'
            });
    }

})();
