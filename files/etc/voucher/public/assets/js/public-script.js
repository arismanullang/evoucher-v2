$( document ).ready(function() {

	var x = findGetParameter('x')+"=";
	console.log(x);
	getProfile(x);
});

function getProfile(x){
	$.ajax({
		url: '/v1/public/redeem/profile?x='+x,
		type: 'get',
		success: function (data) {
			console.log(data);
			$("#holdername").html("");
		}
	});
}

(function() {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.select2').select2();
	}
})();
